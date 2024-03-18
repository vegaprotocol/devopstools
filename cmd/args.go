package cmd

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/secrets/hcvault"
	"github.com/vegaprotocol/devopstools/smartcontracts"
	"github.com/vegaprotocol/devopstools/wallet"

	"go.uber.org/zap"
)

type RootArgs struct {
	Debug  bool
	Logger *zap.Logger

	GitHubToken         string
	FileWithGitHubToken string
	HCVaultURL          string

	hcVaultSecretStore *hcvault.HCVaultSecretStore
}

func (ra *RootArgs) getHCVaultSecretStore() (*hcvault.HCVaultSecretStore, error) {
	if ra.hcVaultSecretStore == nil {
		var err error
		ra.hcVaultSecretStore, err = hcvault.NewHCVaultSecretStore(
			ra.HCVaultURL,
			hcvault.HCVaultLoginToken{
				GitHubToken:         ra.GitHubToken,
				FileWithGitHubToken: ra.FileWithGitHubToken,
			},
		)
		if err != nil {
			return nil, err
		}
	}
	return ra.hcVaultSecretStore, nil
}

func (ra *RootArgs) GetNodeSecretStore() (secrets.NodeSecretStore, error) {
	return ra.getHCVaultSecretStore()
}

func (ra *RootArgs) GetServiceSecretStore() (secrets.ServiceSecretStore, error) {
	return ra.getHCVaultSecretStore()
}

func (ra *RootArgs) GetWalletSecretStore() (secrets.WalletSecretStore, error) {
	return ra.getHCVaultSecretStore()
}

func (ra *RootArgs) GetEthereumClientManager() (*ethutils.EthereumClientManager, error) {
	serviceSecretStore, err := ra.GetServiceSecretStore()
	if err != nil {
		return nil, fmt.Errorf("failed to get EthereumClientManager, %w", err)
	}
	return ethutils.NewEthereumClientManager(serviceSecretStore), nil
}

func (ra *RootArgs) GetSmartContractsManager() (*smartcontracts.SmartContractsManager, error) {
	ethClientManager, err := ra.GetEthereumClientManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get SmartContractsManager, %w", err)
	}
	return smartcontracts.NewSmartContractsManager(ethClientManager, ra.Logger), nil
}

func (ra *RootArgs) GetWalletManager() (*wallet.WalletManager, error) {
	ethClientManager, err := ra.GetEthereumClientManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get SmartContractsManager, %w", err)
	}
	walletSecretStore, err := ra.GetWalletSecretStore()
	if err != nil {
		return nil, fmt.Errorf("failed to get WalletSecretStore, %w", err)
	}
	return wallet.NewWalletManager(ethClientManager, walletSecretStore), nil
}
