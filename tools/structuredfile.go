package tools

import (
	"fmt"

	"github.com/tomwright/dasel"
)

func ReadStructuredFileValue(format, filePath, selector string) (string, error) {
	rootNode, err := dasel.NewFromFile(filePath, format)
	if err != nil {
		return "", fmt.Errorf("failed to open %s file: %w", filePath, err)
	}

	valueNode, err := rootNode.Query(selector)
	if err != nil {
		return "", fmt.Errorf("failed to query given file: %w", err)
	}

	return valueNode.String(), nil
}
