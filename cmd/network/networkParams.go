package network

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/networktools"
	"go.uber.org/zap"
)

type NetworkParamsArgs struct {
	*NetworkArgs
}

var networkParamsArgs NetworkParamsArgs

// networkParamsCmd represents the networkParams command
var networkParamsCmd = &cobra.Command{
	Use:   "network-params",
	Short: "Get Vega Network /network/parameters",
	Long:  `Get Vega Network /network/parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunNetworkParams(networkParamsArgs); err != nil {
			networkParamsArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	networkParamsArgs.NetworkArgs = &networkArgs

	NetworkCmd.AddCommand(networkParamsCmd)
}

func RunNetworkParams(args NetworkParamsArgs) error {
	network, err := networktools.NewNetworkTools(args.VegaNetworkName, args.Logger)
	if err != nil {
		return err
	}
	networkParams, err := network.GetNetworkParams()
	if err != nil {
		return err
	}
	byteNetworkParams, err := json.MarshalIndent(networkParams, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to parse networkParams for network '%s', %w", args.VegaNetworkName, err)
	}
	fmt.Println(string(byteNetworkParams))
	return nil
}
