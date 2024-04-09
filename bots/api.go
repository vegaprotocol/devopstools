package bots

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vegaprotocol/devopstools/tools"

	"go.uber.org/zap"
)

var ErrEmptyResponse = errors.New("empty response")

type (
	ResearchBots map[string]ResearchBot
	TradingBots  map[string]TradingBot
)

func RetrieveResearchBots(ctx context.Context, apiURL, apiKey string, logger *zap.Logger) (ResearchBots, error) {
	leftAttempts := 10
	return tools.RetryReturn(leftAttempts, 5*time.Second, func() (ResearchBots, error) {
		logger.Debug("Retrieving research bots...", zap.String("url", apiURL))

		var payload struct {
			Traders ResearchBots `json:"traders"`
		}
		if err := retrieveBots(ctx, apiURL, apiKey, &payload); err != nil {
			leftAttempts -= 1
			logger.Debug("Failed to retrieve research bots", zap.Error(err), zap.Int("left-attempts", leftAttempts))
			return nil, err
		}

		return payload.Traders, nil
	})
}

func RetrieveTradingBots(ctx context.Context, apiURL, apiKey string, logger *zap.Logger) (TradingBots, error) {
	leftAttempts := 10
	return tools.RetryReturn(leftAttempts, 5*time.Second, func() (TradingBots, error) {
		logger.Debug("Retrieving trading bots...", zap.String("url", apiURL))

		var payload struct {
			Traders TradingBots `json:"traders"`
		}
		if err := retrieveBots(ctx, apiURL, apiKey, &payload); err != nil {
			leftAttempts -= 1
			logger.Debug("Failed to retrieve trading bots", zap.Error(err), zap.Int("left-attempts", leftAttempts))
			return nil, err
		}
		return payload.Traders, nil
	})
}

func retrieveBots(ctx context.Context, apiURL string, apiKey string, payload any) error {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			IdleConnTimeout:       15 * time.Second,
			DisableCompression:    true,
			TLSHandshakeTimeout:   15 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
		},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to build the request: %w", err)
	}
	if len(apiKey) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to submit the request: %w", err)
	}
	defer func() {
		if res.Body != nil {
			_ = res.Body.Close()
		}
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response not ok: %s", res.Status)
	}

	if res.Body == nil {
		return ErrEmptyResponse
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return fmt.Errorf("failed to decode the response payload: %w", err)
	}
	return nil
}
