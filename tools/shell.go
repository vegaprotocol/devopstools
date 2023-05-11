package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func ExecuteBinary(binaryPath string, args []string, structuredOutput interface{}) ([]byte, error) {
	command := exec.Command(binaryPath, args...)

	var stdOut, stErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stErr

	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute binary %s %v with error: %s: %s", binaryPath, args, stErr.String(), err.Error())
	}

	if structuredOutput == nil {
		return stdOut.Bytes(), nil
	}

	if err := json.Unmarshal(stdOut.Bytes(), structuredOutput); err != nil {
		return nil, err
	}

	return nil, nil
}
