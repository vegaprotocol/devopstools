package bots

type TradingBot struct {
	Name       string `json:"name"`
	PubKey     string `json:"pubKey"`
	Parameters struct {
		Base                              string  `json:"marketBase"`
		Quote                             string  `json:"marketQuote"`
		SettlementEthereumContractAddress string  `json:"marketSettlementEthereumContractAddress"`
		SettlementVegaAssetID             string  `json:"marketSettlementVegaAssetID"`
		WantedTokens                      float64 `json:"wantedTokens"`
		CurrentBalance                    float64 `json:"balance"`
		EnableTopUp                       bool    `json:"enableTopUp"`
	} `json:"parameters"`
}
