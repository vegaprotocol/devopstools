package tools

import (
	"fmt"
	"os"
	"strings"
)

func CurrentUserHomePath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return dirname
}

func WhoAmI() (string, error) {
	out, err := ExecuteBinary("whoami", []string{}, nil)
	if err != nil {
		return "", fmt.Errorf("failed to check whoami: %w", err)
	}

	return strings.Trim(string(out), " \t\n"), nil
}
