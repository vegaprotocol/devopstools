package veganetwork

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type VegaNetwork struct {
	Name        string
	DNSSuffix   string
	logger      *zap.Logger
	restTimeout time.Duration
}

func NewVegaNetwork(
	name string,
	logger *zap.Logger,
) (*VegaNetwork, error) {
	var network = VegaNetwork{
		Name:        name,
		logger:      logger,
		restTimeout: time.Millisecond * 500,
	}

	switch name {
	case "devnet":
		network.DNSSuffix = "d.vega.xyz"
	case "stagnet":
		network.DNSSuffix = "s.vega.xyz"
	case "testnet",
		"fairground":
		network.DNSSuffix = "s.vega.xyz"
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
