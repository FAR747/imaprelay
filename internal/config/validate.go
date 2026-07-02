package config

import "fmt"

func (c *Config) SetDefaults() {
	for i := range c.IMAPs {
		if c.IMAPs[i].Security == "" {
			switch c.IMAPs[i].Port {
			case 143:
				c.IMAPs[i].Security = "starttls"
			default:
				c.IMAPs[i].Security = "tls"
			}
		}
	}
}

func (c *Config) Validate() error {
	if c.Targets.Discord == nil && c.Targets.Telegram == nil {
		return fmt.Errorf("at least one target is required")
	}

	if c.Targets.Discord != nil {
		if c.Targets.Discord.WebhookURL == "" {
			return fmt.Errorf("targets.discord.webhook_url is required")
		}
	}

	if c.Targets.Telegram != nil {
		if c.Targets.Telegram.BotToken == "" {
			return fmt.Errorf("targets.telegram.bot_token is required")
		}

		if c.Targets.Telegram.ChatID == "" {
			return fmt.Errorf("targets.telegram.chat_id is required")
		}
	}

	if len(c.IMAPs) == 0 {
		return fmt.Errorf("at least one imap account is required")
	}

	defaultTargets := c.DefaultTargetNames()

	for _, imap := range c.IMAPs {
		if imap.Name == "" {
			return fmt.Errorf("imap.name is required")
		}

		if imap.Host == "" {
			return fmt.Errorf("imap %q: host is required", imap.Name)
		}

		if imap.Port <= 0 {
			return fmt.Errorf("imap %q: port is required and must be greater than 0", imap.Name)
		}

		if imap.Username == "" {
			return fmt.Errorf("imap %q: username is required", imap.Name)
		}

		if imap.Password == "" {
			return fmt.Errorf("imap %q: password is required", imap.Name)
		}

		if imap.Mailbox == "" {
			return fmt.Errorf("imap %q: mailbox is required", imap.Name)
		}

		switch imap.Security {
		case "tls", "starttls", "none":
		default:
			return fmt.Errorf("imap %q: invalid security %q, allowed: tls, starttls, none", imap.Name, imap.Security)
		}

		if imap.Targets == nil {
			if len(defaultTargets) == 0 {
				return fmt.Errorf("imap %q: no targets specified and no default targets configured", imap.Name)
			}
			continue
		}

		if len(*imap.Targets) == 0 {
			return fmt.Errorf("imap %q: targets cannot be empty", imap.Name)
		}

		for _, targetName := range *imap.Targets {
			if !c.HasTarget(targetName) {
				return fmt.Errorf("imap %q: unknown target %q", imap.Name, targetName)
			}
		}
	}

	if c.Proxy != nil {
		if c.Proxy.Type == "" {
			return fmt.Errorf("proxy.type is required")
		}

		switch c.Proxy.Type {
		case "socks5", "http":
			// ok
		default:
			return fmt.Errorf("proxy.type must be socks5 or http")
		}

		if c.Proxy.Address == "" {
			return fmt.Errorf("proxy.address is required")
		}
	}

	return nil
}

func (c *Config) DefaultTargetNames() []string {
	var targets []string

	if c.Targets.Discord != nil && c.Targets.Discord.Default {
		targets = append(targets, "discord")
	}

	if c.Targets.Telegram != nil && c.Targets.Telegram.Default {
		targets = append(targets, "telegram")
	}

	return targets
}

func (c *Config) HasTarget(name string) bool {
	switch name {
	case "discord":
		return c.Targets.Discord != nil
	case "telegram":
		return c.Targets.Telegram != nil
	default:
		return false
	}
}
