package script

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type NetworkActionArgs struct {
	*ScriptArgs
	VegaNetworkName string
}

var networkActionArgs NetworkActionArgs

// networkActionCmd represents the networkAction command
var networkActionCmd = &cobra.Command{
	Use:   "network-action",
	Short: "Quickly do a custom action on running vega network",
	Long:  `Quickly do a custom action on running vega network`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunNetworkAction(networkActionArgs); err != nil {
			networkActionArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	networkActionArgs.ScriptArgs = &scriptArgs

	ScriptCmd.AddCommand(networkActionCmd)

	networkActionCmd.PersistentFlags().StringVar(&networkActionArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := networkActionCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunNetworkAction(args NetworkActionArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	fmt.Printf("Perform some actions on network\n%v\n", network)

	return nil
}
