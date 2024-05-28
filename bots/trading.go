package bots

const (
	// MarketMakerWalletIndex defines index of the market marker wallet. This is hardcoded in the vega-market-sim.
	MarketMakerWalletIndex = 3
)

type BotTraderWantedToken struct {
	PartyId         string  `json:"party_id"`
	Symbol          string  `json:"symbol"`
	VegaAssetId     string  `json:"vega_asset_id"`
	AssetERC20Asset string  `json:"asset_erc20_address"`
	Balance         float64 `json:"balance"`
	WantedTokens    float64 `json:"wanted_tokens"`
}

type TradingBot struct {
	Name       string `json:"name"`
	PubKey     string `json:"pubKey"`
	Parameters struct {
		Base         string                 `json:"marketBase"`
		Quote        string                 `json:"marketQuote"`
		WantedTokens []BotTraderWantedToken `json:"wantedTokens"`

		EnableTopUp bool `json:"enableTopUp"`
	} `json:"parameters"`
}
