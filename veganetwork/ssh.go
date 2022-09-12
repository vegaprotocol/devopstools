package veganetwork

import "github.com/vegaprotocol/devopstools/ssh"

func (network *VegaNetwork) RunCommandOnEveryNode(
	sshUsername string,
	sshPrivateKeyfile string,
	command string,
) map[string]ssh.RunResults {
	return ssh.RunCommandOnEveryHost(network.GetNetworkNodes(), sshUsername, sshPrivateKeyfile, command)
}
