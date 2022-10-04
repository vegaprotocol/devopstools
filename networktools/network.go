package networktools

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type NetworkTools struct {
	Name        string
	DNSSuffix   string
	logger      *zap.Logger
	restTimeout time.Duration
}

func NewNetworkTools(
	name string,
	logger *zap.Logger,
) (*NetworkTools, error) {
	var network = NetworkTools{
		Name:        name,
		logger:      logger,
		restTimeout: time.Second,
	}

	switch name {
	case "devnet":
		network.DNSSuffix = "d.vega.xyz"
	case "fairground":
		network.DNSSuffix = "testnet.vega.xyz"
	case
		"devnet1",
		"devnet2",
		"devnet3",
		"stagnet1",
		"stagnet2",
		"stagnet3":
		network.DNSSuffix = fmt.Sprintf("%s.vega.xyz", name)
	case "mainnet":
		network.DNSSuffix = ""
	default:
		return nil, fmt.Errorf("Unknown network %s", name)
	}
	return &network, nil
}
