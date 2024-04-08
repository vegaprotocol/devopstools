package networktools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vegaprotocol/devopstools/types"

	"go.uber.org/zap"
)

func (network *NetworkTools) GetNetworkParams() (*types.NetworkParams, error) {
	if network.networkParams != nil {
		return network.networkParams, nil
	}
	healthyNodes := network.GetNetworkHealthyNodes()
	for _, node := range healthyNodes {
		params, err := network.GetNetworkParamsFromHost(node, false)
		if err == nil {
			network.networkParams = params
			return params, nil
		}
	}
	return nil, fmt.Errorf("failed to get network parameters for %s network, there is %d healthy nodes", network.Name, len(healthyNodes))
}

func (network *NetworkTools) GetNetworkParamsFromHost(host string, tlsOnly bool) (*types.NetworkParams, error) {
	statsURLs := []string{
		fmt.Sprintf("https://%s/network/parameters", host),
	}
	if !tlsOnly {
		statsURLs = append(statsURLs, fmt.Sprintf("http://%s:3003/network/parameters", host))
		statsURLs = append(statsURLs, fmt.Sprintf("http://%s:3009/network/parameters", host))
	}
	httpClient := http.Client{
		Timeout: network.restTimeout,
	}
	for _, statsURL := range statsURLs {
		req, err := http.NewRequest(http.MethodGet, statsURL, nil)
		if err != nil {
			network.logger.Debug("failed to create new request", zap.String("url", statsURL), zap.Error(err))
			continue
		}
		res, err := httpClient.Do(req)
		if err != nil {
			network.logger.Debug("failed to send request", zap.String("url", statsURL), zap.Error(err))
			continue
		}
		if res.Body == nil {
			network.logger.Debug("response body is empty for request", zap.String("url", statsURL))
			continue
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			network.logger.Debug("failed to read response body", zap.String("url", statsURL), zap.Error(err))
			_ = res.Body.Close()
			continue
		}
		_ = res.Body.Close()

		params := struct {
			NetworkParameters []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"networkParameters"`
		}{}

		if err = json.Unmarshal(body, &params); err != nil {
			network.logger.Debug("failed to parse json for response body", zap.String("url", statsURL), zap.Error(err))
			continue
		}
		result := map[string]string{}
		for _, p := range params.NetworkParameters {
			result[p.Key] = p.Value
		}
		network.logger.Debug("network parameters", zap.String("node", host), zap.Any("parameters", params))
		return types.NewNetworkParams(result), nil
	}
	return nil, fmt.Errorf("failed to get network parameters for host %s", host)
}
