package main

import (
	"flag"
	"fmt"
	"github.com/FAR747/imaprelay/internal/app"
	"github.com/FAR747/imaprelay/internal/config"
	"os"
)

var version = "dev" // VERSION (-ldflags "-X main.version=v0.1.0")

func main() {
	configPath := flag.String("config", "./config.yaml", "path to config file")
	checkConfig := flag.Bool("check-config", false, "validate config and exit")
	flag.Parse()

	cfg, err := config.Load(*configPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	if *checkConfig {
		fmt.Println("config OK")
		fmt.Printf("imaps: %d\n", len(cfg.IMAPs))
		fmt.Printf("imap names: %v\n", cfg.IMAPNames())
		fmt.Printf("enabled targets: %v\n", cfg.EnabledTargetNames())
		fmt.Printf("default targets: %v\n", cfg.DefaultTargetNames())
		return
	}

	fmt.Printf("Starting ImapRelay %s...\n", version)
	fmt.Printf("config: %s\n", *configPath)
	fmt.Printf("imaps: %d\n", len(cfg.IMAPs))

	if err := app.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "app error: %v\n", err)
		os.Exit(1)
	}
}
