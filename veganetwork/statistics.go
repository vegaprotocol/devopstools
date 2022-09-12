package veganetwork

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

type Statistics struct {
	Status      string `json:"status"`
	BlockHeight uint64 `json:"blockHeight,string"`
	CurrentTime string `json:"currentTime"`
	VegaTime    string `json:"vegaTime"`
	AppVersion  string `json:"appVersion"`
}

func (network *VegaNetwork) GetRunningStatistics() (*Statistics, error) {
	hosts := network.GetNetworkNodes()
	httpClient := http.Client{
		Timeout: network.restTimeout,
	}
	resultsChannel := make(chan Statistics, len(hosts)*3)
	var wg sync.WaitGroup
	for _, host := range hosts {
		urls := []string{
			fmt.Sprintf("https://%s/statistics", host),
			fmt.Sprintf("http://%s:3003/statistics", host),
			fmt.Sprintf("http://%s:3009/statistics", host),
		}
		for _, url := range urls {
			wg.Add(1)
			go func(statsURL string) {
				defer wg.Done()

				req, err := http.NewRequest(http.MethodGet, statsURL, nil)
				if err != nil {
					network.logger.Debug("failed to create new request", zap.String("url", statsURL), zap.Error(err))
					return
				}
				res, err := httpClient.Do(req)
				if err != nil {
					network.logger.Debug("failed to send request", zap.String("url", statsURL), zap.Error(err))
					return
				}
				if res.Body == nil {
					network.logger.Debug("response body is empty for request", zap.String("url", statsURL))
					return
				}
				defer res.Body.Close()

				body, err := io.ReadAll(res.Body)
				if err != nil {
					network.logger.Debug("failed to read response body", zap.String("url", statsURL), zap.Error(err))
					return
				}

				stats := struct {
					Statistics Statistics `json:"statistics"`
				}{}

				if err = json.Unmarshal(body, &stats); err != nil {
					network.logger.Debug("failed to parse json for response body", zap.String("url", statsURL), zap.Error(err))
					return
				}
				resultsChannel <- stats.Statistics
			}(url)
		}
	}
	wg.Wait()
	close(resultsChannel)
	var result *Statistics = nil
	for singleResult := range resultsChannel {
		if result == nil || result.BlockHeight < singleResult.BlockHeight {
			result = &singleResult
		}
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("failed to get statistics for network %s, network might be down.", network.Name)
}

func (network *VegaNetwork) GetRunningVersion() (string, error) {
	stats, err := network.GetRunningStatistics()
	if err != nil {
		return "", err
	}
	return stats.AppVersion, nil
}
