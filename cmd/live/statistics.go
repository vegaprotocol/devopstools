package live

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"go.uber.org/zap"
)

type StatisticsArgs struct {
	LiveArgs
	Version     bool
	BlockHeight bool
	All         bool
}

var statisticsArgs StatisticsArgs

// statisticsCmd represents the statistics command
var statisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "Get Vega Network /statistics",
	Long:  `Get Vega Network /statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		statisticsArgs.LiveArgs = liveArgs
		if err := RunStatistics(
			statisticsArgs,
			rootCmd.Logger,
		); err != nil {
			rootCmd.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	LiveCmd.AddCommand(statisticsCmd)
	statisticsCmd.PersistentFlags().BoolVar(&statisticsArgs.All, "all", false, "Get statistics for all hosts")
	statisticsCmd.PersistentFlags().BoolVar(&statisticsArgs.Version, "version", false, "Print version only")
	statisticsCmd.PersistentFlags().BoolVar(&statisticsArgs.BlockHeight, "block", false, "Print block height only")
}

func RunStatistics(
	args StatisticsArgs,
	logger *zap.Logger,
) error {
	network, err := veganetwork.NewVegaNetwork(args.VegaNetworkName, logger)
	if err != nil {
		return err
	}
	if args.All {
		statistics := network.GetRunningStatisticsForAllHosts()
		byteStatistics, err := json.MarshalIndent(statistics, "", "\t")
		if err != nil {
			return fmt.Errorf("failed to parse statistics for network '%s', %w", args.VegaNetworkName, err)
		}
		fmt.Println(string(byteStatistics))

	} else {
		statistics, err := network.GetRunningStatistics()
		if err != nil {
			return err
		}
		if args.Version {
			fmt.Println(statistics.AppVersion)
		} else if args.BlockHeight {
			fmt.Println(statistics.BlockHeight)
		} else {
			byteStatistics, err := json.MarshalIndent(statistics, "", "\t")
			if err != nil {
				return fmt.Errorf("failed to parse statistics for network '%s', %w", args.VegaNetworkName, err)
			}
			fmt.Println(string(byteStatistics))
		}
	}
	return nil
}
