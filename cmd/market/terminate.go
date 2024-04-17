package market

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

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
	*Args

	AllMarkets     bool
	ManagedMarkets bool
	MarketIds      []string
}

var terminateArgs TerminateArgs

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
	terminateArgs.Args = &args

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
	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Debug("Network file loaded", zap.String("name", cfg.Name.String()))

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Debug("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Debug("Looking for healthy gRPC endpoints...")
	healthyEndpoints := tools.FilterHealthyGRPCEndpoints(endpoints)
	if len(healthyEndpoints) == 0 {
		return fmt.Errorf("no healthy gRPC endpoint found on configured datanodes")
	}
	logger.Debug("Healthy gRPC endpoints found", zap.Strings("endpoints", healthyEndpoints))

	datanodeClient := datanode.New(healthyEndpoints, 3*time.Second, args.Logger.Named("datanode"))

	logger.Debug("Connecting to a datanode's gRPC endpoint...")
	dialCtx, cancelDialing := context.WithTimeout(ctx, 2*time.Second)
	defer cancelDialing()
	datanodeClient.MustDialConnection(dialCtx) // blocking
	logger.Debug("Connected to a datanode's gRPC node", zap.String("node", datanodeClient.Target()))

	logger.Debug("Retrieving network parameters...")
	networkParams, err := datanodeClient.GetAllNetworkParameters()
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}
	logger.Debug("Network parameters retrieved")

	minClose, err := time.ParseDuration(networkParams.Params[netparams.GovernanceProposalMarketMinClose])
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	minEnact, err := time.ParseDuration(networkParams.Params[netparams.GovernanceProposalMarketMinEnact])
	if err != nil {
		return fmt.Errorf("failed to parse duration for %q: %w", netparams.GovernanceProposalMarketMinEnact, err)
	}

	marketsToRemove, err := findMarkets(datanodeClient, args.AllMarkets, args.ManagedMarkets, args.MarketIds)
	if err != nil {
		return fmt.Errorf("failed to find markets to remove: %w", err)
	}

	whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	if err := vega.GenerateKeysUpToKey(whaleWallet, cfg.Network.Wallets.VegaTokenWhale.PublicKey); err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	whalePublicKey := cfg.Network.Wallets.VegaTokenWhale.PublicKey

	networkParametersToUpdate := networkParametersForMarketsTermination(networkParams.Params, len(marketsToRemove))
	if len(networkParametersToUpdate) > 0 {
		logger.Info("Voting network parameters required for markets termination", zap.Any("network-parameters", networkParametersToUpdate))

		if _, err := governance.ProposeAndVoteOnNetworkParameters(ctx, networkParametersToUpdate, whaleWallet, whalePublicKey, networkParams, datanodeClient, logger); err != nil {
			return fmt.Errorf("failed to update network parameters required for market termination: %w", err)
		}

		time.Sleep(5 * time.Second)
		logger.Info("Network parameters updated")
	}

	for _, marketDetails := range marketsToRemove {
		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)

		logger.Info("Terminating market", zap.String("market", marketDetails.name))
		proposal := governance.TerminateMarketProposal(closingTime, enactmentTime, marketDetails.name, marketDetails.id, "10")

		args.Logger.Info("Terminating market: Sending proposal", zap.String("market", marketDetails.name))
		proposalRequest := walletpb.SubmitTransactionRequest{
			Command: &walletpb.SubmitTransactionRequest_ProposalSubmission{
				ProposalSubmission: proposal,
			},
		}

		if _, err := walletpkg.SendTransaction(ctx, whaleWallet, whalePublicKey, &proposalRequest, datanodeClient); err != nil {
			return err
		}

		logger.Info("Terminating market: Waiting for proposal to be picked up on the network", zap.String("market", marketDetails.name))
		proposalId, err := tools.RetryReturn(6, 10*time.Second, func() (string, error) {
			reference := proposal.Reference

			res, err := datanodeClient.ListGovernanceData(ctx, &v2.ListGovernanceDataRequest{
				ProposalReference: &reference,
			})
			if err != nil {
				return "", fmt.Errorf("failed to list governance data for reference %s: %w", proposal.Reference, err)
			}
			var proposalId string
			for _, edge := range res.Connection.Edges {
				if edge.Node.Proposal.Reference == reference {
					logger.Info("Found proposal", zap.String("market", marketDetails.name), zap.String("reference", reference),
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
		logger.Info("Terminating market: Proposal found", zap.String("market", marketDetails.name), zap.String("proposal-id", proposalId))

		args.Logger.Info("Terminating market: Voting on proposal", zap.String("market", marketDetails.name), zap.String("proposal-id", proposalId))
		voteRequest := walletpb.SubmitTransactionRequest{
			Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
				VoteSubmission: &commandspb.VoteSubmission{
					ProposalId: proposalId,
					Value:      vegapb.Vote_VALUE_YES,
				},
			},
		}

		if _, err := walletpkg.SendTransaction(ctx, whaleWallet, whalePublicKey, &voteRequest, datanodeClient); err != nil {
			return err
		}

		logger.Info("Terminating market: Voted on proposal", zap.String("market", marketDetails.name), zap.String("proposal-id", proposalId))
	}

	logger.Info("Markets to remove",
		zap.Int("total-markets", len(marketsToRemove)),
		zap.Any("markets", marketsToRemove),
	)
	fmt.Println()

	return nil
}
