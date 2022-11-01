package networktools

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vegaprotocol/devopstools/types"
)

type Traders struct {
	ByERC20TokenHexAddress map[string][]string
	ByFakeAssetId          map[string][]string
}

//
// TRADERBOT
//

func (network *NetworkTools) GetTraderbotBaseURL() (string, error) {
	switch network.Name {
	case "devnet":
		return "https://traderbot-devnet-k8s.ops.vega.xyz", nil
	case types.NetworkStagnet3:
		return "https://traderbot-stagnet3-k8s.ops.vega.xyz", nil
	case types.NetworkFairground:
		return "https://traderbot-testnet-k8s.ops.vega.xyz", nil
	default:
		return fmt.Sprintf("https://traderbot-%s-k8s.ops.vega.xyz", network.Name), nil
	}
}

type TraderbotResponse struct {
	Traders map[string]struct {
		PubKey     string `json:"pubKey"`
		Parameters struct {
			// MarketBase                              string `json:"marketBase"`
			// MarketQuote                             string `json:"marketQuote"`
			SettlementERC20TokenAddress string `json:"marketSettlementEthereumContractAddress"`
			SettlementVegaAssetID       string `json:"marketSettlementVegaAssetID"`
		} `json:"parameters"`
	} `json:"traders"`
}

func (network *NetworkTools) GetTraderbotTraders() (*Traders, error) {
	errMsg := fmt.Sprintf("failed to get traderbot traders for %s", network.Name)
	baseURL, err := network.GetTraderbotBaseURL()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	url := fmt.Sprintf("%s/traders", baseURL)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 5,
	}

	res, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	defer res.Body.Close()

	var payload TraderbotResponse
	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}

	result := Traders{
		ByERC20TokenHexAddress: map[string][]string{},
		ByFakeAssetId:          map[string][]string{},
	}

	for _, trader := range payload.Traders {
		tokenHexAddress := trader.Parameters.SettlementERC20TokenAddress
		if len(tokenHexAddress) > 0 {
			_, ok := result.ByERC20TokenHexAddress[tokenHexAddress]
			if ok {
				result.ByERC20TokenHexAddress[tokenHexAddress] = append(result.ByERC20TokenHexAddress[tokenHexAddress], trader.PubKey)
			} else {
				result.ByERC20TokenHexAddress[tokenHexAddress] = []string{trader.PubKey}
			}
		} else {
			assetId := trader.Parameters.SettlementVegaAssetID
			_, ok := result.ByFakeAssetId[assetId]
			if ok {
				result.ByFakeAssetId[assetId] = append(result.ByFakeAssetId[assetId], trader.PubKey)
			} else {
				result.ByFakeAssetId[assetId] = []string{trader.PubKey}
			}
		}
	}
	return &result, nil
}

//
// LIQBOT
//

func (network *NetworkTools) GetLiqbotBaseURL() (string, error) {
	switch network.Name {
	case "devnet":
		return "https://liqbot-devnet-k8s.ops.vega.xyz", nil
	case types.NetworkStagnet3:
		return "https://liqbot-stagnet3-k8s.ops.vega.xyz", nil
	case types.NetworkFairground:
		return "https://liqbot-testnet-k8s.ops.vega.xyz", nil
	default:
		return fmt.Sprintf("https://liqbot-%s-k8s.ops.vega.xyz", network.Name), nil
	}
}

type liqbotResponse []struct {
	Name                        string `json:"name"`
	PubKey                      string `json:"pubKey"`
	SettlementERC20TokenAddress string `json:"settlementEthereumContractAddress"`
	SettlementVegaAssetID       string `json:"settlementVegaAssetID"`
}

func (network *NetworkTools) GetLiqbotTraders() (*Traders, error) {
	errMsg := fmt.Sprintf("failed to get liqbot traders for %s", network.Name)
	baseURL, err := network.GetLiqbotBaseURL()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	url := fmt.Sprintf("%s/traders-settlement", baseURL)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 5,
	}
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	defer res.Body.Close()

	var payload liqbotResponse
	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}

	result := Traders{
		ByERC20TokenHexAddress: map[string][]string{},
		ByFakeAssetId:          map[string][]string{},
	}

	for _, trader := range payload {
		tokenHexAddress := trader.SettlementERC20TokenAddress
		if len(tokenHexAddress) > 0 {
			_, ok := result.ByERC20TokenHexAddress[tokenHexAddress]
			if ok {
				result.ByERC20TokenHexAddress[tokenHexAddress] = append(result.ByERC20TokenHexAddress[tokenHexAddress], trader.PubKey)
			} else {
				result.ByERC20TokenHexAddress[tokenHexAddress] = []string{trader.PubKey}
			}
		} else {
			assetId := trader.SettlementVegaAssetID
			_, ok := result.ByFakeAssetId[assetId]
			if ok {
				result.ByFakeAssetId[assetId] = append(result.ByFakeAssetId[assetId], trader.PubKey)
			} else {
				result.ByFakeAssetId[assetId] = []string{trader.PubKey}
			}
		}
	}
	return &result, nil
}
