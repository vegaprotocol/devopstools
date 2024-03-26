package bots

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func getBots(ctx context.Context, apiURL string, apiKey string, payload any) error {
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
