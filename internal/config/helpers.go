package config

// Get IMAP names
func (c *Config) IMAPNames() []string {
	names := make([]string, 0, len(c.IMAPs))

	for _, imap := range c.IMAPs {
		names = append(names, imap.Name)
	}

	return names
}

// Get Targets names
func (c *Config) EnabledTargetNames() []string {
	var targets []string

	if c.Targets.Discord != nil {
		targets = append(targets, "discord")
	}

	if c.Targets.Telegram != nil {
		targets = append(targets, "telegram")
	}

	return targets
}
