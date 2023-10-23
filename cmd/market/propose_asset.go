package market

import (
	"fmt"
	"os"
	"time"

	"code.vegaprotocol.io/vega/core/netparams"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
)

type ProposeAssetArgs struct {
	*MarketArgs
}

var proposeAssetArgs ProposeAssetArgs

// provideLPCmd represents the provideLP command
var proposeAssetCmd = &cobra.Command{
	Use:   "propose-asset",
	Short: "Propose asset",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ProposeAssetRun(&proposeAssetArgs); err != nil {
			terminateArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	proposeAssetArgs.MarketArgs = &marketArgs
	MarketCmd.AddCommand(proposeAssetCmd)
}

func ProposeAssetRun(args *ProposeAssetArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	minClose, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinClose])
	if err != nil {
		return fmt.Errorf("failed to parse duration for %s: %w", netparams.GovernanceProposalMarketMinClose, err)
	}
	minEnact, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinEnact])
	if err != nil {
		return fmt.Errorf("failed to parse duration for %s: %w", netparams.GovernanceProposalMarketMinEnact, err)
	}

	closingTime := time.Now().Add(time.Second * 20).Add(minClose)
	enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)

	args.Logger.Info("Proposing asset")
	proposal := governance.NewAssetProposal(closingTime, enactmentTime)

	args.Logger.Info("Proposing asset: Sending proposal")
	walletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: network.VegaTokenWhale.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_ProposalSubmission{
			ProposalSubmission: proposal,
		},
	}
	if err := governance.SubmitTx("propose asset", network.DataNodeClient, network.VegaTokenWhale, args.Logger, &walletTxReq); err != nil {
		return err
	}

	args.Logger.Info("Proposing asset: Waiting for proposal to be picked up on the network")
	proposalId, err := tools.RetryReturn(6, 10*time.Second, func() (string, error) {
		reference := proposal.Reference

		res, err := network.DataNodeClient.ListGovernanceData(&v2.ListGovernanceDataRequest{
			ProposalReference: &reference,
		})
		if err != nil {
			return "", fmt.Errorf("failed to list governance data for reference %s: %w", proposal.Reference, err)
		}
		var proposalId string
		for _, edge := range res.Connection.Edges {
			if edge.Node.Proposal.Reference == reference {
				args.Logger.Info("Found proposal", zap.String("reference", reference),
					zap.String("status", edge.Node.Proposal.State.String()),
					zap.Any("proposal", edge.Node.Proposal))
				proposalId = edge.Node.Proposal.Id
				break
			}
		}

		if len(proposalId) < 1 {
			return "", fmt.Errorf("got empty proposal id for the %s reference", reference)
		}

		return proposalId, nil
	})

	if err != nil {
		return fmt.Errorf("failed to find proposal for new asset:%w", err)
	}
	args.Logger.Info("Proposing asset: Proposal found", zap.String("proposal-id", proposalId))

	args.Logger.Info("Proposing asset: Voting on proposal", zap.String("proposal-id", proposalId))
	voteWalletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: network.VegaTokenWhale.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
			VoteSubmission: &commandspb.VoteSubmission{
				ProposalId: proposalId,
				Value:      vega.Vote_VALUE_YES,
			},
		},
	}
	if err := governance.SubmitTx("vote on asset proposal", network.DataNodeClient, network.VegaTokenWhale, args.Logger, &voteWalletTxReq); err != nil {
		return err
	}

	args.Logger.Info("Proposing asset: Voted on proposal", zap.String("proposal-id", proposalId))

	return nil
}
