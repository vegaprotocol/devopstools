package market

import (
	"fmt"
	"os"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"go.uber.org/zap"
)

type ProposeArgs struct {
	*MarketArgs
}

var proposeArgs ProposeArgs

// proposeCmd represents the propose command
var proposeCmd = &cobra.Command{
	Use:   "propose",
	Short: "Propose and vote on market",
	Long:  `Propose and vote on market`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPropose(proposeArgs); err != nil {
			proposeArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	proposeArgs.MarketArgs = &marketArgs

	MarketCmd.AddCommand(proposeCmd)
}

func runPropose(args ProposeArgs) error {
	// markets := governance.MainnetUpgradeBatchProposal(closeTime, enactTime)
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return fmt.Errorf("failed to create vega network manager: %w", err)
	}

	allMarkets, err := network.DataNodeClient.GetAllMarkets()
	if err != nil {
		return fmt.Errorf("failed to get all markets from the data-node api: %w", err)
	}
	missingMarketsProposals := filterMarkets(args.Logger, allMarkets, MarketProposalsForEnvironment(network.Network))

	defer network.Disconnect()
	closeTime := time.Now().Add(30 * time.Second)
	// enactTime := closeTime.Add(60 * time.Second)

	marketsBatchProposal := governance.NewBatchProposal(
		fmt.Sprintf("%s devopstools markets proposal", network.Network),
		"Create all markets managed by devopstools",
		closeTime,
		missingMarketsProposals,
		nil,
	)

	return sendBatchProposal(args.Logger, network, marketsBatchProposal)
}

func filterMarkets(logger *zap.Logger, allMarkets []*vega.Market, allNetworkProposals []*commandspb.ProposalSubmission) []*commandspb.ProposalSubmission {
	result := []*commandspb.ProposalSubmission{}

	for idx, proposal := range allNetworkProposals {
		newMarketProposal, ok := proposal.Terms.Change.(*vega.ProposalTerms_NewMarket)
		if !ok {
			continue
		}

		// TODO: check if fields are not nill
		marketCode := newMarketProposal.NewMarket.Changes.Instrument.Code
		found := false
		for _, market := range allMarkets {
			// Market is already existing on the network
			if marketCode == market.TradableInstrument.Instrument.Code {
				found = true
				break
			}
		}

		if !found {
			result = append(result, allNetworkProposals[idx])
		}
	}

	return result
}

func sendBatchProposal(logger *zap.Logger, network *veganetwork.VegaNetwork, proposals *commandspb.BatchProposalSubmission) error {
	proposerVegawallet := network.VegaTokenWhale

	// Prepare vegawallet Transaction Request
	walletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_BatchProposalSubmission{
			BatchProposalSubmission: proposals,
		},
	}

	proposalId, err := governance.SubmitTxWithSignature("BatchProposal", network.DataNodeClient, proposerVegawallet, logger, &walletTxReq)
	if err != nil {
		return fmt.Errorf("failed to submit batch proposal with signature: %w", err)
	}

	if err = tools.RetryRun(10, 6*time.Second, func() error {
		return governance.VoteOnProposal("BatchProposal vote", proposalId, proposerVegawallet, network.DataNodeClient, logger)
	}); err != nil {
		return fmt.Errorf("failed to vote on batch proposal(%s): %w", proposalId, err)
	}

	return nil
}
