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

	// New Fairground
	TEURO_F VegaAsset = VegaAsset{
		Name:       "tEURO TEST",
		HexAddress: "0x0158031158Bb4dF2AD02eAA31e8963E84EA978a4",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}
	TDAI_F VegaAsset = VegaAsset{
		Name:       "tDAI TEST",
		HexAddress: "0x26223f9C67871CFcEa329975f7BC0C9cB8FBDb9b",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}
	TUSDC_F VegaAsset = VegaAsset{
		Name:       "tUSDC TEST",
		HexAddress: "0xdBa6373d0DAAAA44bfAd663Ff93B1bF34cE054E9",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}
	TBTC_F VegaAsset = VegaAsset{
		Name:       "tBTC TEST",
		HexAddress: "0x1d525fB145Af5c51766a89706C09fE07E6058D1D",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}

	// OLD Deprecated to be removed
	TEURORopsten VegaAsset = VegaAsset{
		Name:       "tEURO",
		HexAddress: "0xD52b6C949E35A6E4C64b987B1B192A8608931a7b",
		EthNetwork: types.ETHRopsten,
		Version:    erc20token.ERC20TokenOld,
	}
	TDAIRopsten VegaAsset = VegaAsset{
		Name:       "tDAI TEST",
		HexAddress: "0xF4A2bcC43D24D14C4189Ef45fCf681E870675333",
		EthNetwork: types.ETHRopsten,
		Version:    erc20token.ERC20TokenOld,
	}
	TUSDCRopsten VegaAsset = VegaAsset{
		Name:       "tUSDC",
		HexAddress: "0x3773A5c7aFF77e014cBF067dd31801b4C6dc4136",
		EthNetwork: types.ETHRopsten,
		Version:    erc20token.ERC20TokenOld,
	}
	TBTCRopsten VegaAsset = VegaAsset{
		Name:       "tBTC",
		HexAddress: "0xC912F059b4eCCEF6C969B2E0e2544A1A2581C094",
		EthNetwork: types.ETHRopsten,
		Version:    erc20token.ERC20TokenOld,
	}
)

var allAssets = []VegaAsset{TEURO, TDAI, TUSDC, TBTC, TEURO_F, TDAI_F, TUSDC_F, TBTC_F, TEURORopsten, TDAIRopsten, TUSDCRopsten, TBTCRopsten}

func (m *SmartContractsManager) GetAssetWithName(name string) (*erc20token.ERC20Token, error) {
	name = strings.ToLower(name)
	for _, asset := range allAssets {
		if strings.ToLower(asset.Name) == name {
			return m.GetAsset(asset.HexAddress, asset.EthNetwork, asset.Version)
		}
	}
	return nil, fmt.Errorf("there is not token with name %s", name)
}

func (m *SmartContractsManager) GetAssetWithAddress(hexAddress string) (*erc20token.ERC20Token, error) {
	hexAddress = strings.ToLower(hexAddress)
	for _, asset := range allAssets {
		if strings.ToLower(asset.HexAddress) == hexAddress {
			return m.GetAsset(asset.HexAddress, asset.EthNetwork, asset.Version)
		}
	}
	return nil, fmt.Errorf("there is not token with address %s", hexAddress)
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
