package app

import (
	"context"
	"fmt"
	"github.com/FAR747/imaprelay/internal/config"
	"github.com/FAR747/imaprelay/internal/imapclient"
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
		fmt.Printf(
			"- uid=%d from=%q title=%q\nbody=\n%q\nreceived=%s\n",
			msg.UID,
			msg.From,
			msg.Title,
			msg.Body,
			msg.ReceivedAt.Format("2006-01-02 15:04:05"),
		)

		seenUIDs = append(seenUIDs, msg.UID)
	}

	if err := imapclient.MarkSeen(ctx, account, cfg.Proxy, seenUIDs); err != nil {
		return fmt.Errorf("mark seen: %w", err)
	}
	fmt.Printf("Marked messages as seen: account=%s count=%d\n", account.Name, len(seenUIDs))

	return nil
}
