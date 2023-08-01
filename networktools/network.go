package networktools

import (
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/types"
	"go.uber.org/zap"
)

type NetworkTools struct {
	Name          string
	DNSSuffix     string
	logger        *zap.Logger
	restTimeout   time.Duration
	networkParams *types.NetworkParams
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
	case types.NetworkFairground:
		network.DNSSuffix = "testnet.vega.xyz"
	case types.NetworkMainnet:
		network.DNSSuffix = "vega.community"
	case types.NetworkDevnet1:
		network.DNSSuffix = fmt.Sprintf("%s.vega.rocks", name)
	default:
		network.DNSSuffix = fmt.Sprintf("%s.vega.xyz", name)
	}
	return &network, nil
}
