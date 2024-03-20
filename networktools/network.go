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
	network := NetworkTools{
		Name:        name,
		logger:      logger,
		restTimeout: time.Second,
	}

	switch name {
	case types.NetworkFairground:
		network.DNSSuffix = "testnet.vega.rocks"
	case types.NetworkMainnet:
		network.DNSSuffix = "vega.community"
	default:
		network.DNSSuffix = fmt.Sprintf("%s.vega.rocks", name)
	}
	return &network, nil
}
