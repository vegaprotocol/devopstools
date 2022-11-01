package market

import (
	"github.com/vegaprotocol/devopstools/types"
)

const (
	MARKET_AAPL_MARKER    = "auto:aapl"
	MARKET_AAVEDAI_MARKER = "auto:aavedai"
	MARKET_BTCUSD_MARKER  = "auto:btcusd"
	MARKET_ETHBTC_MARKER  = "auto:ethbtc"
	MARKET_TSLA_MARKER    = "auto:tsla"
	MARKET_UNIDAI_MARKER  = "auto:unidai"
	MARKET_ETHDAI_MARKER  = "auto:ethdai"
)

type networkAssetsIDs struct {
	AAPL    string
	AAVEDAI string
	BTCUSD  string
	ETHBTC  string
	TSLA    string
	UNIDAI  string
	ETHDAI  string
}

var settlementAssetIDs map[string]networkAssetsIDs = map[string]networkAssetsIDs{
	types.NetworkDevnet1: {
		AAPL:    "fUSDC",
		AAVEDAI: "fDAI",
		BTCUSD:  "fDAI",
		ETHBTC:  "fBTC",
		TSLA:    "fEURO",
		UNIDAI:  "fDAI",
		ETHDAI:  "fDAI",
	},
	types.NetworkStagnet1: {
		AAPL:    "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d", // "tUSDC"
		AAVEDAI: "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		BTCUSD:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		ETHBTC:  "cee709223217281d7893b650850ae8ee8a18b7539b5658f9b4cc24de95dd18ad", // "tBTC"
		TSLA:    "177e8f6c25a955bd18475084b99b2b1d37f28f3dec393fab7755a7e69c3d8c3b", // "tEURO"
		UNIDAI:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		ETHDAI:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
	},
	types.NetworkStagnet3: {
		AAPL:    "758cd6c8be67b2fb360a97e94200398b417633aeb534fadcfde7e2e42b491750", // "tUSDC"
		AAVEDAI: "74eeaeb035d444a77599dc7e2829ea7fcdc685fd47173a2a2e5fae0994a9d494", // "tDAI"
		BTCUSD:  "74eeaeb035d444a77599dc7e2829ea7fcdc685fd47173a2a2e5fae0994a9d494", // "tDAI"
		ETHBTC:  "5250f006831ceb27b83ac7068a049d6b4f60cf0f1254ccc3c4c11d8952e8cc1e", // "tBTC"
		TSLA:    "37b6a93ad671e25a47682d8562c3c6f4283d051f229b81eba182bd5ae4afa17c", // "tEURO"
		UNIDAI:  "74eeaeb035d444a77599dc7e2829ea7fcdc685fd47173a2a2e5fae0994a9d494", // "tDAI"
		ETHDAI:  "74eeaeb035d444a77599dc7e2829ea7fcdc685fd47173a2a2e5fae0994a9d494", // "tDAI"
	},
	types.NetworkFairground: {
		AAPL:    "c9fe6fc24fce121b2cc72680543a886055abb560043fda394ba5376203b7527d", // "tUSDC"
		AAVEDAI: "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		BTCUSD:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
		ETHBTC:  "cee709223217281d7893b650850ae8ee8a18b7539b5658f9b4cc24de95dd18ad", // "tBTC"
		TSLA:    "177e8f6c25a955bd18475084b99b2b1d37f28f3dec393fab7755a7e69c3d8c3b", // "tEURO"
		UNIDAI:  "b340c130096819428a62e5df407fd6abe66e444b89ad64f670beb98621c9c663", // "tDAI"
	},
}
