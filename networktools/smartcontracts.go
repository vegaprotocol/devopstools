package networktools

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/veganetworksmartcontracts"
)

func (network *NetworkTools) GetEthNetwork() (types.ETHNetwork, error) {
	errMsg := fmt.Sprintf("failed to get Eth Network for %s network", network.Name)
	params, err := network.GetNetworkParams()
	if err != nil {
		return "unknown", fmt.Errorf("%s, %w", errMsg, err)
	}
	ethConfig, err := params.PrimaryEthereumConfig()
	if err != nil {
		return "unknown", fmt.Errorf("%s, %w", errMsg, err)
	}
	ethNetwork, err := types.GetEthNetworkForId(ethConfig.ChainId)
	if err != nil {
		return "unknown", fmt.Errorf("%s, %w", errMsg, err)
	}
	return ethNetwork, nil
}

func (network *NetworkTools) GetSmartContracts(
	ethClientManager *ethutils.EthereumClientManager,
) (*veganetworksmartcontracts.VegaNetworkSmartContracts, error) {
	errMsg := fmt.Sprintf("failed to get Smart Contracts for %s network", network.Name)
	params, err := network.GetNetworkParams()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	ethConfig, err := params.PrimaryEthereumConfig()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	ethNetwork, err := network.GetEthNetwork()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	ethClient, err := ethClientManager.GetEthClient(ethNetwork)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	result, err := veganetworksmartcontracts.NewVegaNetworkSmartContracts(
		ethClient,
		"", // will be taken from Staking Bridge
		"", // will be taken from ERC20 Bridge
		ethConfig.CollateralBridgeContract.Address,
		ethConfig.MultisigControlContract.Address,
		ethConfig.StakingBridgeContract.Address,
		network.logger,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	return result, nil
}
