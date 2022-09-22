package cmd

import (
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/secrets/hcvault"
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
