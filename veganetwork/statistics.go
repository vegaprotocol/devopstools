package veganetwork

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	vegaAPI1 "code.vegaprotocol.io/vega/protos/vega/api/v1"
	"go.uber.org/zap"
)

func (network *VegaNetwork) GetRunningStatistics() (*vegaAPI1.Statistics, error) {
	hosts := network.GetNetworkNodes()
	httpClient := http.Client{
		Timeout: network.restTimeout,
	}
	for _, host := range hosts {
		urls := []string{
			fmt.Sprintf("https://%s/statistics", host),
			fmt.Sprintf("http://%s:3003/statistics", host),
			fmt.Sprintf("http://%s:3009/statistics", host),
		}
		for _, url := range urls {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				network.logger.Debug("failed to create new request", zap.String("url", url), zap.Error(err))
				continue
			}
			res, err := httpClient.Do(req)
			if err != nil {
				network.logger.Debug("failed to send request", zap.String("url", url), zap.Error(err))
				continue
			}
			if res.Body == nil {
				network.logger.Debug("response body is empty for request", zap.String("url", url))
				continue
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				network.logger.Debug("failed to read response body", zap.String("url", url), zap.Error(err))
				continue
			}

			stats := struct {
				Statistics vegaAPI1.Statistics `json:"statistics"`
			}{}

			if err = json.Unmarshal(body, &stats); err != nil {
				network.logger.Debug("failed to parse json for response body", zap.String("url", url), zap.Error(err))
				continue
			}
			// TODO: validate if not too far in history
			return &stats.Statistics, nil
		}
	}
	return nil, fmt.Errorf("failed to get version for network %s, network might be down.", network.Name)
}

func (network *VegaNetwork) GetRunningVersion() (string, error) {
	stats, err := network.GetRunningStatistics()
	if err != nil {
		return "", err
	}
	return stats.AppVersion, nil
}
