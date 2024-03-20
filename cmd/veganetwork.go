package cmd

import (
	"github.com/vegaprotocol/devopstools/networktools"
	"github.com/vegaprotocol/devopstools/veganetwork"
)

func (ra *RootArgs) ConnectToVegaNetwork(network string) (*veganetwork.VegaNetwork, error) {
	tools, err := networktools.NewNetworkTools(network, ra.Logger)
	if err != nil {
		return nil, err
	}
	dataNodeClient, err := tools.GetDataNodeClient()
	if err != nil {
		return nil, err
	}
	nodeSecretStore, err := ra.GetNodeSecretStore()
	if err != nil {
		return nil, err
	}
	serviceSecretStore, err := ra.GetServiceSecretStore()
	if err != nil {
		return nil, err
	}
	primaryEthClientManager, err := ra.GetPrimaryEthereumClientManager()
	if err != nil {
		return nil, err
	}
	primarySmartContractsManager, err := ra.GetPrimarySmartContractsManager()
	if err != nil {
		return nil, err
	}
	secondaryEthClientManager, err := ra.GetSecondaryEthereumClientManager()
	if err != nil {
		return nil, err
	}
	secondarySmartContractsManager, err := ra.GetSecondarySmartContractsManager()
	if err != nil {
		return nil, err
	}
	walletManager, err := ra.GetWalletManager()
	if err != nil {
		return nil, err
	}
	return veganetwork.NewVegaNetwork(
		network,
		dataNodeClient,
		nodeSecretStore,
		serviceSecretStore,
		primaryEthClientManager,
		secondaryEthClientManager,
		primarySmartContractsManager,
		secondarySmartContractsManager,
		walletManager,
		ra.Logger,
	)
}
