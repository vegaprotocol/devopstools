package smartcontracts

import "github.com/vegaprotocol/devopstools/ethutils"

type SmartContractsManager struct {
	ethClientManager *ethutils.EthereumClientManager
	ethURL           string
}

func NewSmartContractsManager(
	ethClientManager *ethutils.EthereumClientManager,
) *SmartContractsManager {
	return &SmartContractsManager{
		ethClientManager: ethClientManager,
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
