package ethutils

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vegaprotocol/devopstools/etherscan"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/types"
)

type EthereumClientManager struct {
	mutex          sync.Mutex
	serviceSecrets secrets.ServiceSecretStore

	ethClientByNetwork map[types.ETHNetwork]*ethclient.Client
}

func NewEthereumClientManager(
	serviceSecrets secrets.ServiceSecretStore,
) *EthereumClientManager {
	return &EthereumClientManager{
		serviceSecrets:     serviceSecrets,
		ethClientByNetwork: map[types.ETHNetwork]*ethclient.Client{},
	}
}

func (m *EthereumClientManager) GetEthereumURL(ethNetwork types.ETHNetwork) (string, error) {
	var ethereumURL string

	switch ethNetwork {

	case types.ETHMainnet, types.ETHSepolia, types.ETHGoerli, types.ETHRopsten:
		if m.serviceSecrets == nil {
			return "", fmt.Errorf("failed to get Ethereum URL for %s, Service Secret Store not provided", ethNetwork)
		}

		var (
			infuraProjectId string
			err             error
		)
		if ethNetwork != types.ETHSepolia {
			infuraProjectId, err = m.serviceSecrets.GetInfuraProjectId()
			if err != nil {
				return "", fmt.Errorf("failed to get Ethereum URL for %s, cannot get Infura Project Id, %w", ethNetwork, err)
			}
		}

		switch ethNetwork {
		case types.ETHMainnet:
			ethereumURL = fmt.Sprintf("https://mainnet.infura.io/v3/%s", infuraProjectId)
		case types.ETHSepolia:
			ethereumURL, err = m.serviceSecrets.GetEthereumNodeURL("sepolia")
			if err != nil {
				return "", fmt.Errorf("failed to get sepolia url from the vault: %w", err)
			}
		case types.ETHGoerli:
			ethereumURL = fmt.Sprintf("https://goerli.infura.io/v3/%s", infuraProjectId)
		case types.ETHRopsten:
			ethereumURL = fmt.Sprintf("https://ropsten.infura.io/v3/%s", infuraProjectId)
		default:
			return "", fmt.Errorf("failed to get Ethereum URL, ethereum network '%s' not supported", ethNetwork)
		}
	default:
		return "", fmt.Errorf("failed to get ethereum client with infura, not supported ethereum network %s", ethNetwork)
	}
	return ethereumURL, nil
}

func (m *EthereumClientManager) GetEthClient(ethNetwork types.ETHNetwork) (*ethclient.Client, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if ethClient, ok := m.ethClientByNetwork[ethNetwork]; ok {
		return ethClient, nil
	}

	url, err := m.GetEthereumURL(ethNetwork)
	if err != nil {
		return nil, fmt.Errorf("failed to get Ethereum Client, %w", err)
	}
	ethClient, err := ethclient.DialContext(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("failed to get Ethereum Client, %w", err)
	}
	m.ethClientByNetwork[ethNetwork] = ethClient

	return ethClient, nil
}

func (m *EthereumClientManager) GetEtherscanClient(ethNetwork types.ETHNetwork) (*etherscan.EtherscanClient, error) {
	switch ethNetwork {
	case types.ETHMainnet, types.ETHSepolia, types.ETHGoerli, types.ETHRopsten:
		if m.serviceSecrets == nil {
			return nil, fmt.Errorf("failed to get Etherscan Client for %s, Service Secret Store not provided", ethNetwork)
		}
		etherscanApikey, err := m.serviceSecrets.GetEtherscanApikey()
		if err != nil {
			return nil, fmt.Errorf("failed to get Etherscan Client for %s, cannot get Etherscan Apikey, %w", ethNetwork, err)
		}
		etherscanClient, err := etherscan.NewEtherscanClient(ethNetwork, etherscanApikey)
		if err != nil {
			return nil, fmt.Errorf("failed to get Etherscan Client for %s, %w", ethNetwork, err)
		}

		return etherscanClient, nil

	default:
		return nil, fmt.Errorf("failed to get etherscan client, not supported network '%s'", ethNetwork)
	}
}
