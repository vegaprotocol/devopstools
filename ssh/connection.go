package ssh

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func GetSSHConnection(
	serverHost string,
	sshUsername string,
	sshPrivateKeyfile string,
) (*ssh.Client, error) {
	key, err := os.ReadFile(sshPrivateKeyfile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: sshUsername,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 2,
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", serverHost), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ssh server: %w", err)
	}

	return client, nil
}
