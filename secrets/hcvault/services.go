package hcvault

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/types"
)

const (
	serviceVaultRoot = "service"
)

func (c *HCVaultSecretStore) GetInfuraProjectId(bridge types.ETHBridge) (string, error) {
	secret, err := c.GetSecret(serviceVaultRoot, bridge.String(), "infura")
	if err != nil {
		return "", err
	}
	return secret["projectId"].(string), nil
}

func (c *HCVaultSecretStore) GetEthereumNodeURL(bridge types.ETHBridge, environment string) (string, error) {
	// secret, err := c.GetSecret(serviceVaultRoot, bridge.String(), "ethereum-node")
	secret, err := c.GetSecret("service", "primary-bridge", "ethereum-node")
	if err != nil {
		return "", fmt.Errorf("failed to get ethereum node url from the vault: %w", err)
	}

	if _, ok := secret[environment]; !ok {
		return "", fmt.Errorf(
			"secret for the ethereum rpc is missing under the %s/%s/%s vault secret",
			serviceVaultRoot,
			bridge.String(),
			"ethereum-node",
		)
	}

	return secret[environment].(string), nil
}

func (c *HCVaultSecretStore) GetEtherscanApikey(bridge types.ETHBridge) (string, error) {
	secret, err := c.GetSecret(serviceVaultRoot, bridge.String(), "etherscan")
	if err != nil {
		return "", err
	}
	return secret["apikey"].(string), nil
}

// primary-bridge/ethereum-node
// service/primary-bridge/ethereum-node
func (c *HCVaultSecretStore) GetDigitalOceanApiToken() (string, error) {
	secret, err := c.GetSecret(serviceVaultRoot, "digitalocean")
	if err != nil {
		return "", err
	}
	return secret["api_token"].(string), nil
}

func (c *HCVaultSecretStore) GetBotsApiToken() (string, error) {
	secret, err := c.GetSecret(serviceVaultRoot, "bots")
	if err != nil {
		return "", err
	}
	return secret["api_token"].(string), nil
}
