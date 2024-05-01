package ssh

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func RunCommandWithClient(
	client *ssh.Client,
	command string,
) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to run command: %w", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("failed to execute command(%v) on ssh session: %w, stderr: %s", command, err, session.Stderr)
	}

	return b.String(), err
}
