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
	case "fairground":
		network.DNSSuffix = "testnet.vega.xyz"
	case "mainnet":
		network.DNSSuffix = ""
	default:
		network.DNSSuffix = fmt.Sprintf("%s.vega.xyz", name)
	}
	return &network, nil
}
