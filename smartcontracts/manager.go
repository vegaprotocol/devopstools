package smartcontracts

import (
	"github.com/vegaprotocol/devopstools/ethutils"
	"go.uber.org/zap"
)

type SmartContractsManager struct {
	ethClientManager *ethutils.EthereumClientManager
	ethURL           string
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

//
// Instance created with this function provides limited functionality
//

func NewSmartContractsManagerWithEthURL(
	ethURL string,
) *SmartContractsManager {
	return &SmartContractsManager{
		ethURL: ethURL,
	}
}
