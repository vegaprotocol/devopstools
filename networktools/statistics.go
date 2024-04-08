package networktools

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

func (network *NetworkTools) GetRunningStatistics() (*Statistics, error) {
	allHostsStats := network.GetRunningStatisticsForAllHosts()
	var result *Statistics
	for hostname, stats := range allHostsStats {
		if result == nil || result.BlockHeight < stats.BlockHeight {
			el := allHostsStats[hostname]
			result = &el
		}
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("failed to get statistics for network %s, network might be down", network.Name)
}

func (network *NetworkTools) GetRunningStatisticsForAllHosts() map[string]Statistics {
	return network.GetRunningStatisticsForHosts(
		network.GetNetworkNodes(false), false,
	)
}

func (network *NetworkTools) GetRunningStatisticsForAllDataNodes() map[string]Statistics {
	return network.GetRunningStatisticsForHosts(
		network.GetNetworkDataNodes(false), true,
	)
}

func (network *NetworkTools) GetRunningStatisticsForHosts(hosts []string, tlsOnly bool) map[string]Statistics {
	type hostStats struct {
		Host       string
		Statistics *Statistics
		Error      error
	}

	resultsChannel := make(chan hostStats, len(hosts))
	var wg sync.WaitGroup
	for _, host := range hosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			stats, err := network.GetRunningStatisticsForHost(host, tlsOnly)
			resultsChannel <- hostStats{
				Host:       host,
				Statistics: stats,
				Error:      err,
			}
		}(host)
	}
	wg.Wait()
	close(resultsChannel)
	result := map[string]Statistics{}
	for singleResult := range resultsChannel {
		if singleResult.Statistics != nil {
			result[singleResult.Host] = *singleResult.Statistics
		}
	}

	network.logger.Debug("Found node statistics", zap.Int("node-num", len(result)))

	return result
}

func (network *NetworkTools) GetRunningStatisticsForHost(host string, tlsOnly bool) (*Statistics, error) {
	statsURLs := []string{
		fmt.Sprintf("https://%s/statistics", host),
	}
	if !tlsOnly {
		statsURLs = append(statsURLs, fmt.Sprintf("http://%s:3003/statistics", host))
		statsURLs = append(statsURLs, fmt.Sprintf("http://%s:3009/statistics", host))
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

		stats := struct {
			Statistics Statistics `json:"statistics"`
		}{}

		if err = json.Unmarshal(body, &stats); err != nil {
			network.logger.Debug("failed to parse json for response body", zap.String("url", statsURL), zap.Error(err))
			continue
		}
		if len(stats.Statistics.Status) == 0 {
			network.logger.Debug("response is missing required data", zap.String("url", statsURL), zap.String("body", string(body)))
			continue
		}
		network.logger.Debug("stats", zap.String("node", host), zap.Any("statistics", stats))
		return &stats.Statistics, nil
	}
	return nil, fmt.Errorf("failed to get statistics for host %s", host)
}

func (network *NetworkTools) GetRunningVersion() (string, error) {
	stats, err := network.GetRunningStatistics()
	if err != nil {
		return "", err
	}
	return stats.AppVersion, nil
}
