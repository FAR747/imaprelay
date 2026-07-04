package app

import (
	"context"
	"fmt"
	"github.com/FAR747/imaprelay/internal/config"
	"github.com/FAR747/imaprelay/internal/imapclient"
	"github.com/FAR747/imaprelay/internal/target"
	"github.com/FAR747/imaprelay/internal/target/discord"
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

	httpClient, err := target.NewHTTPClient(cfg.Proxy)
	if err != nil {
		return fmt.Errorf("create target http client: %w", err)
	}

	successCount := 0
	var errors []error

	for _, targetName := range targetNames {
		switch targetName {
		case "discord":
			if cfg.Targets.Discord == nil {
				errors = append(errors, fmt.Errorf("discord target is not configured"))
				continue
			}

			text, err := target.FormatMessage(msg, target.FormatDiscord)
			if err != nil {
				errors = append(errors, fmt.Errorf("discord format: %w", err))
				continue
			}
			sender := discord.NewWebhookSender(cfg.Targets.Discord.WebhookURL, httpClient)

			if err := sender.Send(ctx, text); err != nil {
				errors = append(errors, fmt.Errorf("discord: %w", err))
				continue
			}

			successCount++

		case "telegram":
			errors = append(errors, fmt.Errorf("telegram: not implemented")) // TODO: Add Telegram
			continue

		default:
			errors = append(errors, fmt.Errorf("%s: unknown target", targetName))
			continue
		}
	}

	if successCount > 0 {
		if len(errors) > 0 {
			fmt.Printf("partial send success: success=%d failed=%d\n", successCount, len(errors))
		}

		return nil
	}

	return fmt.Errorf("all targets failed: %v", errors)
}

func resolveTargetNames(cfg *config.Config, account config.IMAPConfig) []string {
	if account.Targets != nil {
		return *account.Targets
	}

	return cfg.DefaultTargetNames()
}
