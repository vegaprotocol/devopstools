package networktools

import "github.com/vegaprotocol/devopstools/ssh"

func (network *NetworkTools) RunCommandOnEveryNode(
	sshUsername string,
	sshPrivateKeyfile string,
	command string,
) map[string]ssh.RunResults {
	return ssh.RunCommandOnEveryHost(network.GetNetworkNodes(), sshUsername, sshPrivateKeyfile, command)
}
