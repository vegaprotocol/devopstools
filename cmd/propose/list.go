package propose

import (
	"os"

	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ListArgs struct {
	*ProposeArgs
}

var listArgs ListArgs

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list proposals",
	Long:  `list porposals`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunList(listArgs); err != nil {
			listArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	listArgs.ProposeArgs = &proposeArgs

	ProposeCmd.AddCommand(listCmd)
}

func RunList(args ListArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	var (
		logger = args.Logger
	)

	res, err := network.DataNodeClient.ListGovernanceData(&v2.ListGovernanceDataRequest{})
	if err != nil {
		return err
	}

	for _, edge := range res.Connection.Edges {

		errorDetails := ""
		if edge.Node.Proposal.ErrorDetails != nil {
			errorDetails = *edge.Node.Proposal.ErrorDetails
		}

		logger.Info("proposal",
			zap.String("reason", edge.Node.Proposal.Rationale.Description),
			zap.String("state", edge.Node.Proposal.State.String()),
			zap.String("error", errorDetails),
			zap.Any("prop", edge),
		)
	}

	return nil
}
