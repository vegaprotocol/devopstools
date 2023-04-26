package hcvault

import "fmt"

const (
	serviceVaultRoot = "service"
)

func (c *HCVaultSecretStore) GetInfuraProjectId() (string, error) {
	secret, err := c.GetSecret(serviceVaultRoot, "infura")
	if err != nil {
		return "", err
	}
	return secret["projectId"].(string), nil
}

func (c *HCVaultSecretStore) GetEthereumNodeURL(environment string) (string, error) {
	secret, err := c.GetSecret(serviceVaultRoot, "ethereum-node")
	if err != nil {
		return "", fmt.Errorf("failed to get ethereum node url from the vault: %w", err)
	}

	if _, ok := secret[environment]; !ok {
		return "", fmt.Errorf(
			"secret for the %s ethereum node is missing under the %s/%s/%s vault secret",
			environment,
			serviceVaultRoot,
			"ethereum-node",
			environment,
		)
	}

	return secret[environment].(string), nil
}

func (c *HCVaultSecretStore) GetEtherscanApikey() (string, error) {
	secret, err := c.GetSecret(serviceVaultRoot, "etherscan")
	if err != nil {
		return "", err
	}
	return secret["apikey"].(string), nil
}

func (c *HCVaultSecretStore) GetDigitalOceanApiToken() (string, error) {
	secret, err := c.GetSecret(serviceVaultRoot, "digitalocean")
	if err != nil {
		return "", err
	}
	return secret["api_token"].(string), nil
}
