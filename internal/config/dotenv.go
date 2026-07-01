package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("%s:%d: invalid .env line", path, lineNumber)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if key == "" {
			return fmt.Errorf("%s:%d: empty env key", path, lineNumber)
		}

		value = trimEnvValue(value)

		// Real system env has priority over .env.
		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("%s:%d: set env %q: %w", path, lineNumber, key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func trimEnvValue(value string) string {
	value = strings.TrimSpace(value)

	if len(value) >= 2 {
		first := value[0]
		last := value[len(value)-1]

		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			return value[1 : len(value)-1]
		}
	}

	return value
}
