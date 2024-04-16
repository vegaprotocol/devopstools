package bots

import (
	"context"
	"errors"
	"time"

	"github.com/vegaprotocol/devopstools/tools"

	"go.uber.org/zap"
)

const (
	// MarketMakerWalletIndex defines index of the market marker wallet. This is hardcoded in the vega-market-sim.
	MarketMakerWalletIndex = 3
)

var ErrEmptyResponse = errors.New("empty response")

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

type TradingBots map[string]TradingBot

func RetrieveTradingBots(ctx context.Context, apiURL, apiKey string, logger *zap.Logger) (TradingBots, error) {
	leftAttempts := 10
	return tools.RetryReturn(leftAttempts, 5*time.Second, func() (TradingBots, error) {
		logger.Debug("Retrieving trading bots...", zap.String("url", apiURL))

		var payload struct {
			Traders TradingBots `json:"traders"`
		}
		if err := getBots(ctx, apiURL, apiKey, &payload); err != nil {
			leftAttempts -= 1
			logger.Debug("Failed to retrieve trading bots", zap.Error(err), zap.Int("left-attempts", leftAttempts))
			return nil, err
		}
		return payload.Traders, nil
	})
}
