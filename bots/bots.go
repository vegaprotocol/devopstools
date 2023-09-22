package bots

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type BotTraders struct {
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

type ResearchBot struct {
	BotTraders
	WalletData struct {
		Index          int64   `json:"index"`
		PublicKey      string  `json:"publicKey"`
		RecoveryPhrase *string `json:"recoveryPhrase"`
	} `json:"wallet"`
}

func GetBotTraders(
	network string,
) (map[string]BotTraders, error) {
	return GetBotTradersWithURL(network, "")
}

func GetBotTradersWithURL(
	network string,
	botsURL string,
) (map[string]BotTraders, error) {
	if len(botsURL) == 0 {
		botsURL = fmt.Sprintf("https://%s.bots.vega.rocks/traders", network)
	}
	log.Printf("Getting traders from: %s", botsURL)
	var payload struct {
		Traders map[string]BotTraders `json:"traders"`
	}
	err := getBots(botsURL, "", &payload)
	if err != nil {
		return nil, err
	}
	return payload.Traders, nil
}

func GetResearchBots(
	network string,
	botsAPIToken string,
) (map[string]ResearchBot, error) {
	botsURL := fmt.Sprintf("https://%s.bots.vega.rocks/traders", network)
	log.Printf("Getting research bot traders from: %s", botsURL)
	var payload struct {
		Traders map[string]ResearchBot `json:"traders"`
	}
	err := getBots(botsURL, botsAPIToken, &payload)
	if err != nil {
		return nil, err
	}
	return payload.Traders, nil
}

func getBots(botsURL string, botsAPIToken string, payload any) error {
	errMsg := fmt.Sprintf("failed to get bot traders from '%s'", botsURL)

	httpClient := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			IdleConnTimeout:       15 * time.Second,
			DisableCompression:    true,
			TLSHandshakeTimeout:   15 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
		},
	}
	req, err := http.NewRequest(http.MethodGet, botsURL, nil)
	if err != nil {
		return fmt.Errorf("%s, %w", errMsg, err)
	}
	if len(botsAPIToken) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", botsAPIToken))
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s, %w", errMsg, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s, research bot response status code %s", errMsg, res.Status)
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return fmt.Errorf("%s, %w", errMsg, err)
	}
	return nil
}
