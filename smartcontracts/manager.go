package smartcontracts

import (
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/types"

	"go.uber.org/zap"
)

type Manager struct {
	ethClientManager *ethutils.EthereumClientManager
	logger           *zap.Logger
	assets           []VegaAsset
}

func NewManager(ethClientManager *ethutils.EthereumClientManager, bridge types.ETHBridge, logger *zap.Logger) *Manager {
	manager := &Manager{
		ethClientManager: ethClientManager,
		logger:           logger,
	}

	switch bridge {
	case types.PrimaryBridge:
		manager.assets = PrimaryAssets
	case types.SecondaryBridge:
		manager.assets = SecondaryAssets
	}

	return manager
}
