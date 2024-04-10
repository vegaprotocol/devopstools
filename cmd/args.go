package cmd

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/secrets/hcvault"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/wallet"

	"go.uber.org/zap"
)

type RootArgs struct {
	Debug  bool
	Logger *zap.Logger

	GitHubToken         string
	FileWithGitHubToken string
	HCVaultURL          string

	hcVaultSecretStore *hcvault.SecretStore
}

func (ra *RootArgs) getHCVaultSecretStore() (*hcvault.SecretStore, error) {
	if ra.hcVaultSecretStore == nil {
		var err error
		ra.hcVaultSecretStore, err = hcvault.NewHCVaultSecretStore(
			ra.HCVaultURL,
			hcvault.LoginToken{
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

func (ra *RootArgs) GetPrimaryEthereumClientManager() (*ethutils.EthereumClientManager, error) {
	serviceSecretStore, err := ra.GetServiceSecretStore()
	if err != nil {
		return nil, fmt.Errorf("failed to get PrimaryEthereumClientManager, %w", err)
	}
	return ethutils.NewEthereumClientManager(serviceSecretStore, types.PrimaryBridge), nil
}

func (ra *RootArgs) GetWalletManager() (*wallet.Manager, error) {
	ethClientManager, err := ra.GetPrimaryEthereumClientManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get PrimarySmartContractsManager, %w", err)
	}
	walletSecretStore, err := ra.GetWalletSecretStore()
	if err != nil {
		return nil, fmt.Errorf("failed to get WalletSecretStore, %w", err)
	}
	return wallet.NewWalletManager(ethClientManager, walletSecretStore), nil
}
