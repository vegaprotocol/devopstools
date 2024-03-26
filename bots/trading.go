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
