package network

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/networktools"
	"go.uber.org/zap"
)

type GraphArgs struct {
	*NetworkArgs
}

var graphArgs GraphArgs

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Collect Network connection graph",
	Long:  `Collect information about Networks connection graph`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunGraph(graphArgs); err != nil {
			graphArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	graphArgs.NetworkArgs = &networkArgs

	NetworkCmd.AddCommand(graphCmd)
}

func RunGraph(args GraphArgs) error {
	network, err := networktools.NewNetworkTools(args.VegaNetworkName, args.Logger)
	if err != nil {
		return err
	}

	graph, err := network.GetNetworkGraph()
	if err != nil {
		return err
	}

	graph.Print()

	return nil
}
