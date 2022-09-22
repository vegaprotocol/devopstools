package hcvault

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
