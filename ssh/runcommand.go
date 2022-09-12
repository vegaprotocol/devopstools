package ssh

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

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

type RunResults struct {
	Host   string
	Output string
	Err    error
}

func RunCommandOnEveryHost(
	hosts []string,
	sshUsername string,
	sshPrivateKeyfile string,
	command string,
) map[string]RunResults {
	resultsChannel := make(chan RunResults, len(hosts))
	var wg sync.WaitGroup
	for _, host := range hosts {
		wg.Add(1)
		go func(hostName string) {
			defer wg.Done()

			conn, err := GetSSHConnection(hostName, sshUsername, sshPrivateKeyfile)
			if err != nil {
				resultsChannel <- RunResults{Host: hostName, Err: err}
				return
			}
			defer conn.Close()

			stdout, err := RunCommandWithClient(conn, command)
			if err != nil {
				resultsChannel <- RunResults{Host: hostName, Err: err}
				return
			}
			output := strings.TrimSpace(stdout)
			resultsChannel <- RunResults{Host: hostName, Output: output}
		}(host)
	}
	wg.Wait()
	close(resultsChannel)

	// convert to map
	result := map[string]RunResults{}
	for singleResult := range resultsChannel {
		result[singleResult.Host] = singleResult
	}
	return result
}
