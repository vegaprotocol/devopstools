package tools

import (
	"fmt"
	"strings"
)

func WhoAmI() (string, error) {
	out, err := ExecuteBinary("whoami", []string{}, nil)

	if err != nil {
		return "", fmt.Errorf("failed to check whoami: %w", err)
	}

	return strings.Trim(string(out), " \t\n"), nil
}
