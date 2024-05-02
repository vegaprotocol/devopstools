package config

const (
	NetworkDevnet1    NetworkName = "devnet1"
	NetworkStagnet1   NetworkName = "stagnet1"
	NetworkStagnet3   NetworkName = "stagnet3"
	NetworkFairground NetworkName = "fairground"
	NetworkMainnet    NetworkName = "mainnet"
)

type NetworkName string

func (n NetworkName) String() string {
	return string(n)
}
