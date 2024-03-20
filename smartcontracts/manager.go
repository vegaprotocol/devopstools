package smartcontracts

import (
	"github.com/vegaprotocol/devopstools/ethutils"

	"go.uber.org/zap"
)

type SmartContractsManager struct {
	ethClientManager *ethutils.EthereumClientManager
	logger           *zap.Logger
}

func NewSmartContractsManager(
	ethClientManager *ethutils.EthereumClientManager,
	logger *zap.Logger,
) *SmartContractsManager {
	return &SmartContractsManager{
		ethClientManager: ethClientManager,
		logger:           logger,
	}
}
