package hcvault

import (
	"encoding/json"
	"fmt"

	"github.com/vegaprotocol/devopstools/secrets"
)

const (
	walletVaultRoot = "wallet"
)

//
// Ethereum
//

func (c *SecretStore) GetEthereumWallet(secretPath string) (*secrets.EthereumWalletPrivate, error) {
	path := fmt.Sprintf("ethereum/%s", secretPath)
	secretDataByte, err := c.GetSecretAsByte(walletVaultRoot, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get private data for '%s' wallet; %w", path, err)
	}
	var result secrets.EthereumWalletPrivate
	if err = json.Unmarshal(secretDataByte, &result); err != nil {
		return nil, fmt.Errorf("failed to parse private data for '%s' wallet; %w", path, err)
	}
	return &result, nil
}

func (c *SecretStore) StoreEthereumWallet(secretPath string, secretData *secrets.EthereumWalletPrivate) error {
	path := fmt.Sprintf("ethereum/%s", secretPath)
	secretDataByte, err := json.Marshal(secretData)
	if err != nil {
		return fmt.Errorf("failed to parse private data for '%s' wallet; %w", path, err)
	}
	return c.UpsertSecretFromByte(walletVaultRoot, path, secretDataByte)
}

func (c *SecretStore) DoesEthereumWalletExist(secretPath string) (bool, error) {
	path := fmt.Sprintf("ethereum/%s", secretPath)
	return c.DoesExist(walletVaultRoot, path)
}

//
// Vega Wallet
//

func (c *SecretStore) GetVegaWallet(secretPath string) (*secrets.VegaWalletPrivate, error) {
	path := fmt.Sprintf("vegawallet/%s", secretPath)
	secretDataByte, err := c.GetSecretAsByte(walletVaultRoot, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get private data for '%s' wallet; %w", path, err)
	}
	var result secrets.VegaWalletPrivate
	if err = json.Unmarshal(secretDataByte, &result); err != nil {
		return nil, fmt.Errorf("failed to parse private data for '%s' wallet; %w", path, err)
	}
	return &result, nil
}

func (c *SecretStore) StoreVegaWallet(secretPath string, secretData *secrets.VegaWalletPrivate) error {
	path := fmt.Sprintf("vegawallet/%s", secretPath)
	secretDataByte, err := json.Marshal(secretData)
	if err != nil {
		return fmt.Errorf("failed to parse private data for '%s' wallet; %w", path, err)
	}
	return c.UpsertSecretFromByte(walletVaultRoot, path, secretDataByte)
}

func (c *SecretStore) DoesVegaWalletExist(secretPath string) (bool, error) {
	path := fmt.Sprintf("vegawallet/%s", secretPath)
	return c.DoesExist(walletVaultRoot, path)
}
