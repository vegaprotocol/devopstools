package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TendermintRESTClient struct {
	httpClient *http.Client
}

func NewTendermintRESTClient() *TendermintRESTClient {
	return &TendermintRESTClient{
		httpClient: &http.Client{
			Timeout: time.Second * 1,
		},
	}
}

type NodeInfo struct {
	TendermintNodeId string `json:"id"`
	Moniker          string `json:"moniker"`
	ListenAddr       string `json:"listen_addr"`
}

//
// /net_info
//

type NetInfo struct {
	Peers []struct {
		NodeInfo NodeInfo `json:"node_info"`
		RemoteIp string   `json:"remote_ip"`
	} `json:"peers"`
}

func (c *TendermintRESTClient) GetNodeNetInfo(endpoint string) (*NetInfo, error) {
	var (
		errMsg = "failed to get Tendermint Node Net Info, %w"
	)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	u = u.JoinPath("net_info")
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	defer response.Body.Close()

	var payload struct {
		Result NetInfo `json:"result"`
	}
	if err = json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	return &payload.Result, nil
}

//
// /status
//

type Status struct {
	NodeInfo      NodeInfo `json:"node_info"`
	ValidatorInfo struct {
		VotingPower int64 `json:"voting_power,string"`
	} `json:"validator_info"`
}

func (c *TendermintRESTClient) GetNodeStatus(endpoint string) (*Status, error) {
	var (
		errMsg = "failed to get Tendermint Node Status, %w"
	)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	u = u.JoinPath("status")
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	defer response.Body.Close()

	var payload struct {
		Result Status `json:"result"`
	}
	if err = json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	return &payload.Result, nil
}
