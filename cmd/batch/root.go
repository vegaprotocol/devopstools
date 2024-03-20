package batch

import (
	"fmt"
	"log"
	"time"

	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/tools"

	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"

	"github.com/spf13/cobra"
)

type BatchArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var batchArgs BatchArgs

// Root Command for Incentive
var BatchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Send batch proposals to the network",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := sendBatchProposals(batchArgs); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	batchArgs.RootArgs = &rootCmd.Args

	BatchCmd.PersistentFlags().StringVar(&batchArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := BatchCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func sendBatchProposals(args BatchArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()
	closeTime := time.Now().Add(30 * time.Second)
	enactTime := closeTime.Add(60 * time.Second)
	proposals := governance.MainnetUpgradeBatchProposal(closeTime, enactTime)

	proposerVegawallet := network.VegaTokenWhale

	// Prepare vegawallet Transaction Request
	walletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_BatchProposalSubmission{
			BatchProposalSubmission: proposals,
		},
	}

	proposalId, err := governance.SubmitTxWithSignature("BatchProposal", network.DataNodeClient, proposerVegawallet, args.Logger, &walletTxReq)
	if err != nil {
		return err
	}

	if err = tools.RetryRun(10, 6*time.Second, func() error {
		return governance.VoteOnProposal("BatchProposal vote", proposalId, proposerVegawallet, network.DataNodeClient, args.Logger)
	}); err != nil {
		return fmt.Errorf("failed to vote on batch proposal(%s): %w", proposalId, err)
	}

	return nil
}
