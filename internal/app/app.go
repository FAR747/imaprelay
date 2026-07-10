package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FAR747/imaprelay/internal/config"
)

const DefaultCheckTime = 120 // in seconds

func Run(cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	fmt.Println("ImapRelay started")

	if err := runOnce(ctx, cfg); err != nil {
		fmt.Printf("poll error: %v\n", err)
	}

	checkTime := time.Duration(DefaultCheckTime) * time.Second // TODO: add interval in config
	ticker := time.NewTicker(checkTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("App stopped")
			return nil

		case <-ticker.C:
			if err := runOnce(ctx, cfg); err != nil {
				fmt.Printf("poll error: %v\n", err)
			}
		}
	}
}

func runOnce(ctx context.Context, cfg *config.Config) error {
	for _, account := range cfg.IMAPs {
		if err := processAccount(ctx, cfg, account); err != nil {
			fmt.Printf("account error: account=%s error=%v\n", account.Name, err)
			continue
		}
	}

	return nil
}
