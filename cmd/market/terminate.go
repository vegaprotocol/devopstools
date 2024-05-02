package market

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi"

	"code.vegaprotocol.io/vega/core/netparams"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

const OpsManagedMetadata = "managed:vega/ops"

type TerminateArgs struct {
	AllMarkets     bool
	ManagedMarkets bool
	MarketIds      []string

	*Args
}

var terminateArgs TerminateArgs

// provideLPCmd represents the provideLP command
var terminateCmd = &cobra.Command{
	Use:   "terminate",
	Short: "Terminate one or more markets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunTerminate(&terminateArgs); err != nil {
			terminateArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	terminateArgs.Args = &marketArgs

	terminateCmd.PersistentFlags().BoolVar(&terminateArgs.AllMarkets, "all", false, "Terminate all markets")
	terminateCmd.PersistentFlags().BoolVar(&terminateArgs.ManagedMarkets, "managed", false, fmt.Sprintf("Terminate markets managed by ops only(all with %s metadata)", OpsManagedMetadata))
	terminateCmd.PersistentFlags().StringSliceVar(&terminateArgs.MarketIds, "market-ids", []string{}, "Terminate only certain markets")

	Cmd.AddCommand(terminateCmd)
}

type marketDetails struct {
	name string
	id   string
}

func findMarkets(dataNodeClient vegaapi.DataNodeClient, allMarkets bool, managedMarkets bool, marketIds []string) ([]marketDetails, error) {
	var result []marketDetails

	markets, err := dataNodeClient.GetAllMarkets(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get all markets")
	}

	isManaged := func(market *vegapb.Market) bool {
		if market.TradableInstrument.Instrument.Metadata == nil {
			return false
		}

		for _, metadata := range market.TradableInstrument.Instrument.Metadata.Tags {
			if metadata == OpsManagedMetadata {
				return true
			}
		}

		return false
	}

	for _, market := range markets {
		if !slices.Contains(governance.LiveMarketStates, market.State) {
			continue
		}

		if allMarkets || slices.Contains(marketIds, market.Id) || (managedMarkets && isManaged(market)) {
			result = append(result, marketDetails{
				id:   market.Id,
				name: market.TradableInstrument.Instrument.Name,
			})
		}
	}

	return result, nil
}

func networkParametersForMarketsTermination(currentNetworkParameters map[string]string, marketsToTerminate int) map[string]string {
	const spamParametersMultiplier = 5

	result := map[string]string{}
	maxVotes := currentNetworkParameters[netparams.SpamProtectionMaxVotes]
	maxProposals := currentNetworkParameters[netparams.SpamProtectionMaxProposals]

	maxVotesInt := tools.StrToIntOrDefault(maxVotes, 0)
	maxProposalsInt := tools.StrToIntOrDefault(maxProposals, 0)

	if maxVotesInt < marketsToTerminate*spamParametersMultiplier {
		result[netparams.SpamProtectionMaxVotes] = fmt.Sprintf("%d", marketsToTerminate*spamParametersMultiplier)
	}

	if maxProposalsInt < marketsToTerminate*spamParametersMultiplier {
		result[netparams.SpamProtectionMaxProposals] = fmt.Sprintf("%d", marketsToTerminate*spamParametersMultiplier)
	}

	return result
}

func RunTerminate(args *TerminateArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	ctx := context.Background()

	minClose, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinClose])
	if err != nil {
		return fmt.Errorf("failed to parse duration for %s: %w", netparams.GovernanceProposalMarketMinClose, err)
	}
	minEnact, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinEnact])
	if err != nil {
		return fmt.Errorf("failed to parse duration for %s: %w", netparams.GovernanceProposalMarketMinEnact, err)
	}

	marketsToRemove, err := findMarkets(network.DataNodeClient, args.AllMarkets, args.ManagedMarkets, args.MarketIds)
	if err != nil {
		return fmt.Errorf("failed to find markets to remove: %w", err)
	}

	networkParametersToUpdate := networkParametersForMarketsTermination(network.NetworkParams.Params, len(marketsToRemove))
	if len(networkParametersToUpdate) > 0 {
		args.Logger.Sugar().Infof("Voting network parmeters required for markets termination: %v", networkParametersToUpdate)

		if _, err := governance.ProposeAndVoteOnNetworkParameters(ctx, networkParametersToUpdate, network.VegaTokenWhale, "", network.NetworkParams, network.DataNodeClient, args.Logger); err != nil {
			return fmt.Errorf("failed to update network parameters required for market termination: %w", err)
		}

		time.Sleep(5 * time.Second)
		args.Logger.Info("Network parameters updated")
	}

	firstKey := vega.MustFirstKey(network.VegaTokenWhale)

	for _, marketDetails := range marketsToRemove {
		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)

		args.Logger.Info("Terminating market", zap.String("market", marketDetails.name))
		proposal := governance.TerminateMarketProposal(closingTime, enactmentTime, marketDetails.name, marketDetails.id, "10")

		args.Logger.Info("Terminating market: Sending proposal", zap.String("market", marketDetails.name))
		proposalRequest := walletpb.SubmitTransactionRequest{
			Command: &walletpb.SubmitTransactionRequest_ProposalSubmission{
				ProposalSubmission: proposal,
			},
		}

		if _, err := walletpkg.SendTransaction(ctx, network.VegaTokenWhale, firstKey, &proposalRequest, network.DataNodeClient); err != nil {
			return err
		}

		args.Logger.Info("Terminating market: Waiting for proposal to be picked up on the network", zap.String("market", marketDetails.name))
		proposalId, err := tools.RetryReturn(6, 10*time.Second, func() (string, error) {
			reference := proposal.Reference

			res, err := network.DataNodeClient.ListGovernanceData(ctx, &v2.ListGovernanceDataRequest{
				ProposalReference: &reference,
			})
			if err != nil {
				return "", fmt.Errorf("failed to list governance data for reference %s: %w", proposal.Reference, err)
			}
			var proposalId string
			for _, edge := range res.Connection.Edges {
				if edge.Node.Proposal.Reference == reference {
					args.Logger.Info("Found proposal", zap.String("market", marketDetails.name), zap.String("reference", reference),
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
			return fmt.Errorf("failed to find proposal for terminate market %s: %w", marketDetails.name, err)
		}
		args.Logger.Info("Terminating market: Proposal found", zap.String("market", marketDetails.name), zap.String("proposal-id", proposalId))

		args.Logger.Info("Terminating market: Voting on proposal", zap.String("market", marketDetails.name), zap.String("proposal-id", proposalId))
		voteRequest := walletpb.SubmitTransactionRequest{
			Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
				VoteSubmission: &commandspb.VoteSubmission{
					ProposalId: proposalId,
					Value:      vegapb.Vote_VALUE_YES,
				},
			},
		}

		if _, err := walletpkg.SendTransaction(ctx, network.VegaTokenWhale, firstKey, &voteRequest, network.DataNodeClient); err != nil {
			return err
		}

		args.Logger.Info("Terminating market: Voted on proposal", zap.String("market", marketDetails.name), zap.String("proposal-id", proposalId))
	}

	fmt.Printf("%v\n", marketsToRemove)
	fmt.Println(len(marketsToRemove))

	return nil
}
