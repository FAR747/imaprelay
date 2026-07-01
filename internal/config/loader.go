package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

const DefaultConfigPath = "./config.yaml"

func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigPath
	}

	envPath := filepath.Join(filepath.Dir(path), ".env")
	if fileExists(envPath) {
		if err := LoadDotEnv(envPath); err != nil {
			return nil, fmt.Errorf("load .env: %w", err)
		}
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	expanded, err := ExpandEnvStrict(string(raw))
	if err != nil {
		return nil, fmt.Errorf("expand env: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	cfg.SetDefaults()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
