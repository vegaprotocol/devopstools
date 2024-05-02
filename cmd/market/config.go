package market

import (
	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/governance/market"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
)

const CoinBaseOraclePubKey = "0xfCEAdAFab14d46e20144F48824d0C09B1a03F2BC"

type networkAssetsIDs struct {
	AAPL    string
	AAVEDAI string
	BTCUSD  string
	ETHBTC  string
	TSLA    string
	UNIDAI  string
	ETHDAI  string

	SettlementAsset_USDC  string
	MainnetLikeAsset_USDT string
}

var l2Configs = map[config.NetworkName][]*vega.EthereumL2Config{
	config.NetworkDevnet1: {
		&vega.EthereumL2Config{
			NetworkId:     "100",
			ChainId:       "100",
			Name:          "Gnosis Chain",
			Confirmations: 3,
			BlockInterval: 3,
		},
	},
	config.NetworkStagnet1: {
		&vega.EthereumL2Config{
			NetworkId:     "100",
			ChainId:       "100",
			Name:          "Gnosis Chain",
			Confirmations: 3,
			BlockInterval: 3,
		},
		&vega.EthereumL2Config{
			NetworkId:     "1",
			ChainId:       "1",
			Name:          "Ethereum Mainnet",
			Confirmations: 3,
			BlockInterval: 1,
		},
		&vega.EthereumL2Config{
			NetworkId:     "5",
			ChainId:       "5",
			Name:          "Goerli",
			Confirmations: 3,
			BlockInterval: 3,
		},
	},
	config.NetworkStagnet3: {},
	config.NetworkFairground: {
		&vega.EthereumL2Config{
			NetworkId:     "100",
			ChainId:       "100",
			Name:          "Gnosis Chain",
			Confirmations: 3,
			BlockInterval: 3,
		},
		&vega.EthereumL2Config{
			NetworkId:     "42161",
			ChainId:       "42161",
			Name:          "Arbitrum One",
			Confirmations: 3,
			BlockInterval: 50,
		},
		&vega.EthereumL2Config{
			NetworkId:     "5",
			ChainId:       "5",
			Name:          "Goerli",
			Confirmations: 3,
			BlockInterval: 3,
		},
	},
}

var settlementAssetIDs = map[config.NetworkName]networkAssetsIDs{
	config.NetworkDevnet1: {
		AAPL:    "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d", // "tUSDC"
		AAVEDAI: "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		BTCUSD:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		ETHBTC:  "cee709223217281d7893b650850ae8ee8a18b7539b5658f9b4cc24de95dd18ad", // "fBTC"
		TSLA:    "177e8f6c25a955bd18475084b99b2b1d37f28f3dec393fab7755a7e69c3d8c3b", // "fEURO"
		UNIDAI:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		ETHDAI:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"

		SettlementAsset_USDC:  "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d", // "tUSDC"
		MainnetLikeAsset_USDT: "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d",
	},
	config.NetworkStagnet1: {
		AAPL:    "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d", // "tUSDC"
		AAVEDAI: "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		BTCUSD:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		ETHBTC:  "cee709223217281d7893b650850ae8ee8a18b7539b5658f9b4cc24de95dd18ad", // "tBTC"
		TSLA:    "177e8f6c25a955bd18475084b99b2b1d37f28f3dec393fab7755a7e69c3d8c3b", // "tEURO"
		UNIDAI:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		ETHDAI:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"

		SettlementAsset_USDC:  "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d", // "tUSDC"
		MainnetLikeAsset_USDT: "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d",
	},
	config.NetworkStagnet3: {
		AAPL:    "ede4076aef07fd79502d14326c54ab3911558371baaf697a19d077f4f89de399", // "tUSDC"
		AAVEDAI: "16ae5dbb1fd7aa2ddef725703bfe66b3647a4da7b844bfdd04e985756f53d9d6", // "tDAI"
		BTCUSD:  "16ae5dbb1fd7aa2ddef725703bfe66b3647a4da7b844bfdd04e985756f53d9d6", // "tDAI"
		ETHBTC:  "e1cc8e2598d11c4c3ccc4521f0fc06f4b6d940a8607ca38b72bec138600f0525", // "tBTC"
		TSLA:    "4e4e80abff30cab933b8c4ac6befc618372eb76b2cbddc337eff0b4a3a4d25b8", // "tEURO"
		UNIDAI:  "16ae5dbb1fd7aa2ddef725703bfe66b3647a4da7b844bfdd04e985756f53d9d6", // "tDAI"
		ETHDAI:  "16ae5dbb1fd7aa2ddef725703bfe66b3647a4da7b844bfdd04e985756f53d9d6", // "tDAI"

		SettlementAsset_USDC:  "",
		MainnetLikeAsset_USDT: "",
	},
	config.NetworkFairground: {
		AAPL:    "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d", // "tUSDC"
		AAVEDAI: "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		BTCUSD:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		ETHBTC:  "cee709223217281d7893b650850ae8ee8a18b7539b5658f9b4cc24de95dd18ad", // "tBTC"
		TSLA:    "177e8f6c25a955bd18475084b99b2b1d37f28f3dec393fab7755a7e69c3d8c3b", // "tEURO"
		UNIDAI:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		ETHDAI:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"

		SettlementAsset_USDC:  "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d", // "tUSDC"
		MainnetLikeAsset_USDT: "8ba0b10971f0c4747746cd01ff05a53ae75ca91eba1d4d050b527910c983e27e", // Probably should be fetched dynamicly
	},
}

func ProposalsForEnvironment(environment config.NetworkName) []*commandspb.ProposalSubmission {
	switch environment {
	case config.NetworkDevnet1:
		return []*commandspb.ProposalSubmission{
			market.NewMainnetSimulationBitcoinTetherPerpetualWithoutTime(settlementAssetIDs[config.NetworkDevnet1].MainnetLikeAsset_USDT),
			market.NewMainnetSimulationEtherTetherPerpetualWithoutTime(settlementAssetIDs[config.NetworkDevnet1].MainnetLikeAsset_USDT),
			market.NewFutureBTCUSDTWithoutTime(settlementAssetIDs[config.NetworkDevnet1].MainnetLikeAsset_USDT, CoinBaseOraclePubKey),
			market.NewFutureETHUSDTWithoutTime(settlementAssetIDs[config.NetworkDevnet1].MainnetLikeAsset_USDT, CoinBaseOraclePubKey),
			market.NewMainnetSimulationSNXUSDTPerpWithoutTime(settlementAssetIDs[config.NetworkDevnet1].MainnetLikeAsset_USDT),
			market.NewMainnetSimulationLDOUSDTPerpWithoutTime(settlementAssetIDs[config.NetworkDevnet1].MainnetLikeAsset_USDT),
			market.NewMainnetSimulationINJUSDTPerpWithoutTime(settlementAssetIDs[config.NetworkDevnet1].MainnetLikeAsset_USDT),
		}
	case config.NetworkStagnet1:
		return []*commandspb.ProposalSubmission{}
	case config.NetworkStagnet3:
		return []*commandspb.ProposalSubmission{}
	case config.NetworkFairground:
		return []*commandspb.ProposalSubmission{}
	case config.NetworkMainnet:
		return []*commandspb.ProposalSubmission{}
	default:
		return []*commandspb.ProposalSubmission{}
	}
}
