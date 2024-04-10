package network

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vegaprotocol/devopstools/networktools"
	toolslib "github.com/vegaprotocol/devopstools/tools"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type HealthyNodesArgs struct {
	*Args
}

var healthyNodesArgs HealthyNodesArgs

var healthyNodesCmd = &cobra.Command{
	Use:   "healthy-nodes",
	Short: "Get and print healthy nodes from the given network",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunHealthyNodes(healthyNodesArgs); err != nil {
			healthyNodesArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	healthyNodesArgs.Args = &networkArgs

	Cmd.AddCommand(healthyNodesCmd)
}

type output struct {
	Validators          []string `json:"validators"`
	Explorers           []string `json:"explorers"`
	DataNodes           []string `json:"data_nodes"`
	TendermintEndpoints []string `json:"tendermint_endpoints"`
	All                 []string `json:"all"`
}

func RunHealthyNodes(args HealthyNodesArgs) error {
	logger := args.Logger

	logger.Debug("Connecting to the vega network")
	tools, err := networktools.NewNetworkTools(args.VegaNetworkName, logger)
	if err != nil {
		return fmt.Errorf("failed to get network tools: %w", err)
	}

	allNodes := tools.GetNetworkNodes(true)
	blockExplorers := tools.GetBlockExplorers(true)
	dataNodes := tools.GetNetworkDataNodes(true)
	tendermintEndpoints := tools.GetNetworkTendermintRESTEndpoints(true)
	var validators []string

	for _, nodeHost := range allNodes {
		var (
			isDataNode bool
			isExplorer bool
		)
		for _, dataNodeHost := range dataNodes {
			// assuming data node host has the `api.` prefix
			if strings.Contains(dataNodeHost, nodeHost) {
				isDataNode = true
				break
			}
		}

		for _, explorerHost := range blockExplorers {
			// assuming data node host has the `api.` prefix
			if strings.Contains(explorerHost, nodeHost) {
				isExplorer = true
				break
			}
		}

		if !isDataNode && !isExplorer {
			validators = append(validators, nodeHost)
		}
	}

	var healthyValidators []string
	var healthyExplorers []string
	var healthyDataNodes []string

	for _, host := range blockExplorers {
		if err := toolslib.RetryRun(3, 500*time.Millisecond, func() error {
			if !isNodeHealthy(logger, host, false) {
				return fmt.Errorf("node is not healthy")
			}

			return nil
		}); err == nil {
			healthyExplorers = append(healthyExplorers, host)
		}
	}

	for _, host := range validators {
		if err := toolslib.RetryRun(3, 500*time.Millisecond, func() error {
			if !isNodeHealthy(logger, host, false) {
				return fmt.Errorf("node is not healthy")
			}

			return nil
		}); err == nil {
			healthyValidators = append(healthyValidators, host)
		}
	}

	for _, host := range dataNodes {
		if err := toolslib.RetryRun(3, 500*time.Millisecond, func() error {
			if !isNodeHealthy(logger, host, true) {
				return fmt.Errorf("node is not healthy")
			}

			return nil
		}); err == nil {
			healthyDataNodes = append(healthyDataNodes, host)
		}
	}

	allHealthyEndpoints := append(healthyExplorers, append(healthyValidators, healthyDataNodes...)...)

	healthyTendermintNodes := []string{}
	for _, node := range tendermintEndpoints { // e.g: http://api0.vega.community:26657
		for _, healthyEndpoint := range allHealthyEndpoints { // e.g: api2.vega.community
			if strings.Contains(node, healthyEndpoint) {
				healthyTendermintNodes = append(healthyTendermintNodes, node)
				break
			}
		}
	}

	result := output{
		Validators:          healthyValidators,
		Explorers:           healthyExplorers,
		DataNodes:           healthyDataNodes,
		TendermintEndpoints: healthyTendermintNodes,
		All:                 allHealthyEndpoints,
	}

	resp, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	fmt.Println(string(resp))

	return nil
}

type statisticsResponse struct {
	Statistics struct {
		BlockHeight string `json:"blockHeight"`
		CurrentTime string `json:"currentTime"`
		VegaTime    string `json:"vegaTime"`
	} `json:"statistics"`
}

// Simple logic to check data node is healthy and ready to use somewhere, where We want
func isNodeHealthy(logger *zap.Logger, host string, dataNode bool) bool {
	const timeThresholds = 10 * time.Second
	const blocksThresholds = 10

	httpClient := http.Client{
		Timeout: time.Second,
	}
	resp, err := httpClient.Get(fmt.Sprintf("https://%s/statistics", host))
	if err != nil {
		logger.Sugar().Debugf("Failed to GET for node %s: %s", host, err.Error())
		return false
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Sugar().Debugf("Failed to read response body for node %s: %s", host, err.Error())
		return false
	}

	statsResponse := statisticsResponse{}
	if err := json.Unmarshal(responseBytes, &statsResponse); err != nil {
		logger.Sugar().Debugf("Failed to unsmrshal response into golang structure for node %s: %s", host, err.Error())
		return false
	}

	currentTime, err := time.Parse(time.RFC3339Nano, statsResponse.Statistics.CurrentTime)
	if err != nil {
		logger.Sugar().Debugf("Failed to parse current time (%s): %s", statsResponse.Statistics.CurrentTime, err.Error())
		return false
	}

	vegaTime, err := time.Parse(time.RFC3339Nano, statsResponse.Statistics.VegaTime)
	if err != nil {
		logger.Sugar().Debugf("Failed to parse vega time (%s): %s", statsResponse.Statistics.VegaTime, err.Error())
		return false
	}

	if currentTime.Sub(vegaTime) > timeThresholds {
		// Time diff too big
		logger.Sugar().Debugf("Time diff too big")
		return false
	}

	if !dataNode {
		return true // Validator looks healthy.
	}

	dataNodeCurrentBlock := resp.Header.Get("X-Block-Height")
	if dataNodeCurrentBlock == "" {
		logger.Sugar().Debugf("Failed to get X-Block-Height header")
		return false
	}

	vegaBlock, err := strconv.ParseUint(statsResponse.Statistics.BlockHeight, 10, 64)
	if err != nil {
		logger.Sugar().Debugf("failed to convert vega block(%s) to int: %s", statsResponse.Statistics.BlockHeight, err.Error())
		return false
	}

	dataNodeBlock, err := strconv.ParseUint(dataNodeCurrentBlock, 10, 64)
	if err != nil {
		logger.Sugar().Debugf("failed to convert data-node block(%s) to int: %s", dataNodeCurrentBlock, err.Error())
		return false
	}

	if vegaBlock-dataNodeBlock > blocksThresholds {
		return false
	}

	return true
}
