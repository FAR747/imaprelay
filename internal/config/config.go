package config

type Config struct {
	Targets TargetsConfig `yaml:"targets"`
	IMAPs   []IMAPConfig  `yaml:"imaps"`
	Proxy   *ProxyConfig  `yaml:"proxy,omitempty"`
}

type TargetsConfig struct {
	Discord  *DiscordConfig  `yaml:"discord,omitempty"`
	Telegram *TelegramConfig `yaml:"telegram,omitempty"`
}

type DiscordConfig struct {
	WebhookURL string `yaml:"webhook_url"`
	Default    bool   `yaml:"default,omitempty"`
}

type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
	Default  bool   `yaml:"default,omitempty"`
}

type IMAPConfig struct {
	Name     string    `yaml:"name"`
	Host     string    `yaml:"host"`
	Port     int       `yaml:"port"`
	Username string    `yaml:"username"`
	Password string    `yaml:"password"`
	Mailbox  string    `yaml:"mailbox"`
	Security string    `yaml:"security,omitempty"`
	Targets  *[]string `yaml:"targets,omitempty"`
}

type ProxyConfig struct {
	Type    string `yaml:"type"`
	Address string `yaml:"address"`
}
