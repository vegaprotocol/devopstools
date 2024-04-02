package propose

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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
	Short: "get list of proposals",
	Long:  `get list of porposals`,
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

	res, err := network.DataNodeClient.ListGovernanceData(&v2.ListGovernanceDataRequest{})
	if err != nil {
		return err
	}

	for _, edge := range res.Connection.Edges {
		PrettyPrintProposal(edge, args.Debug)
	}

	return nil
}

func PrettyPrintProposal(proposal *v2.GovernanceDataEdge, details bool) {
	fmt.Printf("Proposal: \"%s\" (%s)\n",
		proposal.Node.Proposal.Rationale.Title,
		proposal.Node.Proposal.State.String(),
	)
	if proposal.Node.Proposal.ErrorDetails != nil {
		fmt.Printf("- error: %s\n", *proposal.Node.Proposal.ErrorDetails)
	}
	if details {
		s, _ := json.MarshalIndent(proposal, "\t", "  ")
		fmt.Printf("details: %s\n", string(s))
	} else {
		proposalTime := time.Unix(0, proposal.Node.Proposal.Timestamp)
		enactTime := time.Unix(proposal.Node.Proposal.Terms.EnactmentTimestamp, 0)
		fmt.Printf("- time: %s\n", proposalTime.Format(time.RFC3339))
		fmt.Printf("- enact: %s\n", enactTime.Format(time.RFC3339))
	}
}
