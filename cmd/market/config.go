package market

const (
	MARKET_AAPL_MARKER    = "auto:aapl"
	MARKET_AAVEDAI_MARKER = "auto:aavedai"
	MARKET_BTCUSD_MARKER  = "auto:btcusd"
	MARKET_ETHBTC_MARKER  = "auto:ethbtc"
	MARKET_TSLA_MARKER    = "auto:tsla"
	MARKET_UNIDAI_MARKER  = "auto:unidai"
)

var (
	marketConfig = map[string]struct {
		Marker           string
		FakeAssetSymbol  string
		ERC20AssetSymbol string
	}{
		"AAPL": {
			Marker:           MARKET_AAPL_MARKER,
			FakeAssetSymbol:  "fUSDC",
			ERC20AssetSymbol: "tUSDC",
		},
		"AAVEDAI": {
			Marker:           MARKET_AAVEDAI_MARKER,
			FakeAssetSymbol:  "fDAI",
			ERC20AssetSymbol: "tDAI",
		},
		"BTCUSD": {
			Marker:           MARKET_BTCUSD_MARKER,
			FakeAssetSymbol:  "fDAI",
			ERC20AssetSymbol: "tDAI",
		},
		"ETHBTC": {
			Marker:           MARKET_ETHBTC_MARKER,
			FakeAssetSymbol:  "fBTC",
			ERC20AssetSymbol: "tBTC",
		},
		"TSLA": {
			Marker:           MARKET_TSLA_MARKER,
			FakeAssetSymbol:  "fEURO",
			ERC20AssetSymbol: "tEURO",
		},
		"UNIDAI": {
			Marker:           MARKET_UNIDAI_MARKER,
			FakeAssetSymbol:  "fDAI",
			ERC20AssetSymbol: "tDAI",
		},
	}
)
