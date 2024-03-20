package smartcontracts

import (
	"fmt"
	"strings"

	"github.com/vegaprotocol/devopstools/smartcontracts/erc20token"
	"github.com/vegaprotocol/devopstools/types"
)

type VegaAsset struct {
	Name       string
	HexAddress string
	EthNetwork types.ETHNetwork
	Version    erc20token.ERC20TokenVersion
}

func (m *Manager) GetAssetWithName(name string) (*erc20token.ERC20Token, error) {
	name = strings.ToLower(name)
	for _, asset := range m.assets {
		if strings.ToLower(asset.Name) == name {
			return m.GetAsset(asset.HexAddress, asset.EthNetwork, asset.Version)
		}
	}
	return nil, fmt.Errorf("there is not token with name %s", name)
}

func (m *Manager) GetAssetWithAddress(hexAddress string) (*erc20token.ERC20Token, error) {
	hexAddress = strings.ToLower(hexAddress)
	for _, asset := range m.assets {
		if strings.ToLower(asset.HexAddress) == hexAddress {
			return m.GetAsset(asset.HexAddress, asset.EthNetwork, asset.Version)
		}
	}
	return nil, fmt.Errorf("there is not token with address %s", hexAddress)
}

func (m *Manager) GetAsset(
	hexAddress string,
	ethNetwork types.ETHNetwork,
	version erc20token.ERC20TokenVersion,
) (*erc20token.ERC20Token, error) {
	client, err := m.ethClientManager.GetEthClient(ethNetwork)
	if err != nil {
		return nil, err
	}
	return erc20token.NewERC20Token(client, hexAddress, version)
}
