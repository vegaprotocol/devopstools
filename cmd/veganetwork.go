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
	ethClientManager, err := ra.GetEthereumClientManager()
	if err != nil {
		return nil, err
	}
	smartContractsManager, err := ra.GetSmartContractsManager()
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
		ethClientManager,
		smartContractsManager,
		walletManager,
		ra.Logger,
	)
}
