package secrets

import "github.com/vegaprotocol/devopstools/types"

//
// Secrets for services
//

type ServiceSecretStore interface {
	GetInfuraProjectId(types.ETHBridge) (string, error)
	GetEthereumNodeURL(types.ETHBridge, string) (string, error)
	GetEtherscanApikey(types.ETHBridge) (string, error)
	GetDigitalOceanApiToken() (string, error)
}
