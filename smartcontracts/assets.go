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

var (
	TEURO VegaAsset = VegaAsset{
		Name:       "tEURO",
		HexAddress: "0x7119500C6327928ae4E64531feBB258f0a33A617",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}
	TDAI VegaAsset = VegaAsset{
		Name:       "tDAI",
		HexAddress: "0x973cB2a51F83a707509fe7cBafB9206982E1c3ad",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}
	TUSDC VegaAsset = VegaAsset{
		Name:       "tUSDC",
		HexAddress: "0x40ff2D218740EF033b43B8Ce0342aEBC81934554",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}
	TBTC VegaAsset = VegaAsset{
		Name:       "tBTC",
		HexAddress: "0x123cB4a2AB190F88a50646D5436aB3F5859107Ed",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}
)

func (m *SmartContractsManager) GetAssetWithName(name string) (*erc20token.ERC20Token, error) {
	name = strings.ToLower(name)
	for _, asset := range []VegaAsset{TEURO, TDAI, TUSDC, TBTC} {
		if strings.ToLower(asset.Name) == name {
			return m.GetAsset(asset.HexAddress, asset.EthNetwork, asset.Version)
		}
	}
	return nil, fmt.Errorf("there is not token with name %s", name)
}

func (m *SmartContractsManager) GetAsset(
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
