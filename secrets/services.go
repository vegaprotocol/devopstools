package secrets

//
// Secrets for services
//

type ServiceSecretStore interface {
	GetInfuraProjectId() (string, error)
	GetEthereumNodeURL(string) (string, error)
	GetEtherscanApikey() (string, error)
	GetDigitalOceanApiToken() (string, error)
}
