package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetSourcecodeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		SourceCode           string `json:"SourceCode"`
		ABI                  string `json:"ABI"`
		ContractName         string `json:"ContractName"`
		CompilerVersion      string `json:"CompilerVersion"`
		ConstructorArguments string `json:"ConstructorArguments"`
	} `json:"result"`
}

func (c *Client) GetSourcecode(ctx context.Context, hexAddress string) (*GetSourcecodeResponse, string, error) {
	downloadURL := fmt.Sprintf(
		"%s?module=contract&action=getsourcecode&address=%s",
		c.apiURL, hexAddress,
	)
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, downloadURL, fmt.Errorf("failed rate limiter for Get Sourcecode for smart contract: %s. %w", hexAddress, err)
	}
	downloadURLWithApikey := fmt.Sprintf("%s&apikey=%s", downloadURL, c.apikey)
	resp, err := http.Get(downloadURLWithApikey)
	if err != nil {
		return nil, downloadURL, fmt.Errorf("failed to Get Sourcecode for smart contract: %s. %w", hexAddress, err)
	}
	defer resp.Body.Close()
	var payload GetSourcecodeResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, downloadURL, fmt.Errorf("failed to parse Get Sourcecode response to json: %s. %w", hexAddress, err)
	}
	if payload.Status != "1" {
		return &payload, downloadURL, fmt.Errorf("get Sourcecode response Status is not 1: %s. Error: %s", payload.Status, payload.Message)
	}
	return &payload, downloadURL, nil
}
