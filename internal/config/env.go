package config

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

var envPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

func ExpandEnvStrict(input string) (string, error) {
	missing := make(map[string]struct{})

	output := envPattern.ReplaceAllStringFunc(input, func(match string) string {
		parts := envPattern.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}

		name := parts[1]

		value, ok := os.LookupEnv(name)
		if !ok {
			missing[name] = struct{}{}
			return match
		}

		return value
	})

	if len(missing) > 0 {
		names := make([]string, 0, len(missing))
		for name := range missing {
			names = append(names, name)
		}

		sort.Strings(names)

		return "", fmt.Errorf("missing environment variables: %s", strings.Join(names, ", "))
	}

	return output, nil
}
