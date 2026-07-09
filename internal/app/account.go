package app

import (
	"context"
	"fmt"
	"github.com/FAR747/imaprelay/internal/config"
	"github.com/FAR747/imaprelay/internal/imapclient"
	"github.com/FAR747/imaprelay/internal/target"
	"github.com/FAR747/imaprelay/internal/target/discord"
	"net/http"
)

func processAccount(ctx context.Context, cfg *config.Config, account config.IMAPConfig) error {
	fmt.Printf("Checking IMAP account: %s\n", account.Name)

	messages, err := imapclient.FetchUnseen(ctx, account, cfg.Proxy)
	seenUIDs := make([]imapclient.UID, 0, len(messages))
	if err != nil {
		return fmt.Errorf("fetch unseen: %w", err)
	}

	if len(messages) == 0 {
		fmt.Printf("No unseen messages: account=%s\n", account.Name)
		return nil
	}

	fmt.Printf("Unseen messages: account=%s count=%d\n", account.Name, len(messages))

	for _, msg := range messages {
		//fmt.Println(target.FormatMessage(msg))

		if err := sendMessage(ctx, cfg, account, msg); err != nil {
			fmt.Printf("send error: account=%s uid=%d error=%v\n", account.Name, msg.UID, err)
			continue
		}

		seenUIDs = append(seenUIDs, msg.UID)
	}

	if err := imapclient.MarkSeen(ctx, account, cfg.Proxy, seenUIDs); err != nil {
		return fmt.Errorf("mark seen: %w", err)
	}
	fmt.Printf("Marked messages as seen: account=%s count=%d\n", account.Name, len(seenUIDs))

	return nil
}

func sendMessage(ctx context.Context, cfg *config.Config, account config.IMAPConfig, msg imapclient.Message) error {
	targetNames := resolveTargetNames(cfg, account)

	if len(targetNames) == 0 {
		return fmt.Errorf("no targets resolved")
	}

	httpClient, err := target.NewHTTPClient(cfg.Proxy)
	if err != nil {
		return fmt.Errorf("create target http client: %w", err)
	}

	successCount := 0
	var errors []error

	for _, targetName := range targetNames {
		sender, formatType, err := buildSender(cfg, targetName, httpClient)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		text, err := target.FormatMessage(msg, formatType)
		if err != nil {
			errors = append(errors, fmt.Errorf("%s format: %w", targetName, err))
			continue
		}

		if err := sender.Send(ctx, text); err != nil {
			errors = append(errors, fmt.Errorf("%s send: %w", targetName, err))
			continue
		}

		successCount++
	}

	if successCount > 0 {
		if len(errors) > 0 {
			fmt.Printf("partial send success: success=%d failed=%d\n", successCount, len(errors))
		}

		return nil
	}

	return fmt.Errorf("all targets failed: %v", errors)
}

func buildSender(cfg *config.Config, targetName string, httpClient *http.Client) (target.Sender, target.FormatType, error) {
	switch targetName {
	case "discord":
		if cfg.Targets.Discord == nil {
			return nil, "", fmt.Errorf("discord target is not configured")
		}

		return discord.NewWebhookSender(cfg.Targets.Discord.WebhookURL, httpClient), target.FormatDiscord, nil

	case "telegram":
		return nil, "", fmt.Errorf("telegram target is not implemented yet")

	default:
		return nil, "", fmt.Errorf("unknown target %q", targetName)
	}
}

func resolveTargetNames(cfg *config.Config, account config.IMAPConfig) []string {
	if account.Targets != nil {
		return *account.Targets
	}

	return cfg.DefaultTargetNames()
}
