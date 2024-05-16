package network

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vegaprotocol/devopstools/config"
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
	healthyNodesArgs.Args = &args

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
	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Debug("Network file loaded", zap.String("name", cfg.Name.String()))

	explorersEndpoints := []string{cfg.Explorer.RESTURL}
	datanodesEndpoints := config.ListDatanodeRESTEndpoints(cfg)
	blockchainEndpoints := config.ListBlockchainRESTEndpoints(cfg)
	validatorsEndpoints := config.ListValidatorRESTEndpoints(cfg)

	var healthyValidators []string
	var healthyExplorers []string
	var healthyDataNodes []string

	for _, host := range explorersEndpoints {
		if err := toolslib.RetryRun(3, 500*time.Millisecond, func() error {
			if !isNodeHealthy(ctx, logger, host, false) {
				return fmt.Errorf("node is not healthy")
			}

			return nil
		}); err == nil {
			healthyExplorers = append(healthyExplorers, host)
		}
	}

	for _, host := range validatorsEndpoints {
		if err := toolslib.RetryRun(3, 500*time.Millisecond, func() error {
			if !isNodeHealthy(ctx, logger, host, false) {
				return fmt.Errorf("node is not healthy")
			}

			return nil
		}); err == nil {
			healthyValidators = append(healthyValidators, host)
		}
	}

	for _, host := range datanodesEndpoints {
		if err := toolslib.RetryRun(3, 500*time.Millisecond, func() error {
			if !isNodeHealthy(ctx, logger, host, true) {
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
func isNodeHealthy(ctx context.Context, logger *zap.Logger, host string, dataNode bool) bool {
	const timeThresholds = 10 * time.Second
	const blocksThresholds = 10

	reqCtx, cancelReq := context.WithTimeout(ctx, time.Second)
	defer cancelReq()

	request, err := http.NewRequestWithContext(reqCtx, http.MethodGet, fmt.Sprintf("%s/statistics", host), nil)
	if err != nil {
		logger.Debug("Building query failed", zap.String("host", host), zap.Error(err))
		return false
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Debug("Querying node failed", zap.String("host", host), zap.Error(err))
		return false
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Debug("Failed to read node response", zap.String("host", host), zap.Error(err))
		return false
	}

	statsResponse := statisticsResponse{}
	if err := json.Unmarshal(responseBytes, &statsResponse); err != nil {
		logger.Debug("Failed to unmarshal response", zap.String("host", host), zap.Error(err))
		return false
	}

	currentTime, err := time.Parse(time.RFC3339Nano, statsResponse.Statistics.CurrentTime)
	if err != nil {
		logger.Debug("Failed to parse current time", zap.String("current-time", statsResponse.Statistics.CurrentTime), zap.Error(err))
		return false
	}

	vegaTime, err := time.Parse(time.RFC3339Nano, statsResponse.Statistics.VegaTime)
	if err != nil {
		logger.Debug("Failed to parse vega time", zap.String("current-time", statsResponse.Statistics.VegaTime), zap.Error(err))
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
		logger.Debug("Failed to get X-Block-Height header")
		return false
	}

	vegaBlock, err := strconv.ParseUint(statsResponse.Statistics.BlockHeight, 10, 64)
	if err != nil {
		logger.Debug("Failed to convert vega block", zap.String("block-height", statsResponse.Statistics.BlockHeight), zap.Error(err))
		return false
	}

	dataNodeBlock, err := strconv.ParseUint(dataNodeCurrentBlock, 10, 64)
	if err != nil {
		logger.Debug("Failed to convert data-node block", zap.String("block-height", dataNodeCurrentBlock), zap.Error(err))
		return false
	}

	if vegaBlock-dataNodeBlock > blocksThresholds {
		return false
	}

	return true
}
