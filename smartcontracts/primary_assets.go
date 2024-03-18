package smartcontracts

import (
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20token"
	"github.com/vegaprotocol/devopstools/types"
)

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

	// Sandbox
	VT_S VegaAsset = VegaAsset{
		Name:       "Vega Token Sandbox",
		HexAddress: "0x51d9dbe9a724c6a8383016fad566e55c95359d36",
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

	// New Stagnet3
	WOZ_S3 VegaAsset = VegaAsset{
		Name:       "Woz Token (Vega)",
		HexAddress: "0x559cc3042F28dbaBE30A6b2343c102faeA08D399",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}

	TIM_S3 VegaAsset = VegaAsset{
		Name:       "Tim Token (Vega)",
		HexAddress: "0x1e071110c83876dc71fBB9B85f273e1DDC805F12",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}

	MAK_S3 VegaAsset = VegaAsset{
		Name:       "Maker Reward Token (Vega)",
		HexAddress: "0x4B8cC8de9Dae629dDB3e64A6b4669077AD9aA0C4",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}

	TEURO_S3 VegaAsset = VegaAsset{
		Name:       "tEURO TEST",
		HexAddress: "0xF47c3A0f61ED18386db1FD87Aad3C4523Ec326E8",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}

	TAK_S3 VegaAsset = VegaAsset{
		Name:       "Taker Reward Token (Vega)",
		HexAddress: "0x7698D6c27326eB53a09eF50A6d851e7692cC82da",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}

	TBTC_S3 VegaAsset = VegaAsset{
		Name:       "tBTC TEST",
		HexAddress: "0x333a2B77fd3c261DfAbB8E161d9063F6c15A3816",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}

	TUSDC_S3 VegaAsset = VegaAsset{
		Name:       "tUSDC TEST",
		HexAddress: "0x6b3D260116d9a87458E44718b3DE7fABa8ac745C",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}

	TDAI_S3 VegaAsset = VegaAsset{
		Name:       "tDAI TEST",
		HexAddress: "0x355C3914Ea8F25559D5b8c3E1134c57fB3739B7A",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}

	LIQ_S3 VegaAsset = VegaAsset{
		Name:       "Liquidity Reward Token (Vega)",
		HexAddress: "0x3303C7BcF0aa1858D4c3cE7E372dd10809aF7f86",
		EthNetwork: types.ETHSepolia,
		Version:    erc20token.ERC20TokenBase,
	}
	VEGA_S3 VegaAsset = VegaAsset{
		Name:       "Vega (stagnet3)",
		HexAddress: "0xF136d9Ca8f9C2F6501487994e498fCDC48813Ae6",
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

var PrimaryAssets = []VegaAsset{
	TEURO, TDAI, TUSDC, TBTC, VT_S, TEURO_F, TDAI_F, TUSDC_F, TBTC_F, TEURORopsten, TDAIRopsten, TUSDCRopsten, TBTCRopsten,
	VEGA_S3, LIQ_S3, TDAI_S3, TBTC_S3, TUSDC_S3, TAK_S3, TEURO_S3, MAK_S3, TIM_S3, WOZ_S3,
}
