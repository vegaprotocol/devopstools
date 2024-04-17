package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetTxlistResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		Input string `json:"input"`
	} `json:"result"`
}

func (c *Client) GetTxlist(ctx context.Context, hexAddress string) (*GetTxlistResponse, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("failed rate limiter for Get Txlist for smart contract: %s. %w", hexAddress, err)
	}
	resp, err := http.Get(fmt.Sprintf(
		"%s?apikey=%s&module=account&action=txlist&address=%s",
		c.apiURL, c.apikey, hexAddress,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to get Creation Code for smart contract: %s. %w", hexAddress, err)
	}
	defer resp.Body.Close()
	var payload GetTxlistResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to convert Creation Code response to json: %s. %w", hexAddress, err)
	}
	if payload.Status != "1" {
		return &payload, fmt.Errorf("get Sourcecode response Status is not 1: %s. Error: %s", payload.Status, payload.Message)
	}
	return &payload, nil
}
