package network

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"go.uber.org/zap"
)

type StatsArgs struct {
	*NetworkArgs
	Version     bool
	BlockHeight bool
}

var statsArgs StatsArgs

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get Vega Network /statistics",
	Long:  `Get Vega Network /statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunStats(statsArgs); err != nil {
			statsArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	statsArgs.NetworkArgs = &networkArgs

	NetworkCmd.AddCommand(statsCmd)
	statsCmd.PersistentFlags().BoolVar(&statsArgs.Version, "version", false, "Print version only")
	statsCmd.PersistentFlags().BoolVar(&statsArgs.BlockHeight, "block", false, "Print Block Height only")
}

func RunStats(args StatsArgs) error {
	network, err := veganetwork.NewVegaNetwork(args.VegaNetworkName, args.Logger)
	if err != nil {
		return err
	}
	stats, err := network.GetRunningStatistics()
	if err != nil {
		return err
	}
	if args.Version {
		fmt.Println(stats.AppVersion)
	} else if args.BlockHeight {
		fmt.Println(stats.BlockHeight)
	} else {
		byteStats, err := json.MarshalIndent(stats, "", "\t")
		if err != nil {
			return fmt.Errorf("failed to parse stats for network '%s', %w", args.VegaNetworkName, err)
		}
		fmt.Println(string(byteStats))
	}
	return nil
}
