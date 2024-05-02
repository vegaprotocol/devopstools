package market

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/governance/networkparameters"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/veganetwork"

	"code.vegaprotocol.io/vega/core/netparams"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ProposeArgs struct {
	*Args
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
	proposeArgs.Args = &marketArgs

	Cmd.AddCommand(proposeCmd)
}

func runPropose(args ProposeArgs) error {
	// markets := governance.MainnetUpgradeBatchProposal(closeTime, enactTime)
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return fmt.Errorf("failed to create vega network manager: %w", err)
	}

	allMarkets, err := network.DataNodeClient.GetAllMarkets(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get all markets from the data-node api: %w", err)
	}
	missingMarketsProposals := filterMarkets(args.Logger, allMarkets, ProposalsForEnvironment(network.Network))

	defer network.Disconnect()

	networkParamsProposals, err := preMarketDeployProposals(network.Network, network.NetworkParams)
	if err != nil {
		return fmt.Errorf("failed to determine network parameters to update before markets are proposed: %w", err)
	}

	if len(missingMarketsProposals) > 0 && len(networkParamsProposals) > 0 {
		closeTime := time.Now().Add(30 * time.Second)
		networkParamsBatchProposal := governance.NewBatchProposal(
			fmt.Sprintf("%s devopstools network params proposal", network.Network),
			"Update network parameters before markets are proposed",
			closeTime,
			networkParamsProposals,
			nil,
		)
		args.Logger.Info("Updating network parameters")

		if err := sendBatchProposal(network, networkParamsBatchProposal); err != nil {
			return fmt.Errorf("failed to send batch proposal to update network parameters: %w", err)
		}
	} else {
		args.Logger.Info("No network parameters need to be updated")
	}

	closeTime := time.Now().Add(30 * time.Second)
	// enactTime := closeTime.Add(60 * time.Second)

	marketsBatchProposal := governance.NewBatchProposal(
		fmt.Sprintf("%s devopstools markets proposal", network.Network),
		"Create all markets managed by devopstools",
		closeTime,
		missingMarketsProposals,
		nil,
	)

	if len(missingMarketsProposals) < 1 {
		args.Logger.Info("All required markets exist")
		return nil
	}

	return sendBatchProposal(network, marketsBatchProposal)
}

func filterMarkets(logger *zap.Logger, allMarkets []*vegapb.Market, allNetworkProposals []*commandspb.ProposalSubmission) []*commandspb.ProposalSubmission {
	var result []*commandspb.ProposalSubmission

	for idx, proposal := range allNetworkProposals {
		newMarketProposal, ok := proposal.Terms.Change.(*vegapb.ProposalTerms_NewMarket)
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
			logger.Sugar().Infof("Market %s will be created. Adding proposal to batch", marketCode)
		} else {
			logger.Sugar().Infof("Market %s found on the network. No need to recreate it", marketCode)
		}
	}

	return result
}

func sendBatchProposal(network *veganetwork.VegaNetwork, proposals *commandspb.BatchProposalSubmission) error {
	ctx := context.Background()
	whaleWallet := network.VegaTokenWhale
	whalePublicKey := vega.MustFirstKey(whaleWallet)

	// Prepare vegawallet Transaction Request
	walletTxReq := walletpb.SubmitTransactionRequest{
		Command: &walletpb.SubmitTransactionRequest_BatchProposalSubmission{
			BatchProposalSubmission: proposals,
		},
	}

	resp, err := walletpkg.SendTransaction(ctx, whaleWallet, whalePublicKey, &walletTxReq, network.DataNodeClient)
	if err != nil {
		return fmt.Errorf("failed to submit batch proposal with signature: %w", err)
	}

	proposalId := resp.TxHash

	if err = tools.RetryRun(10, 6*time.Second, func() error {
		return governance.VoteOnProposal(ctx, "BatchProposal vote", proposalId, whaleWallet, whalePublicKey, network.DataNodeClient)
	}); err != nil {
		return fmt.Errorf("failed to vote on batch proposal(%s): %w", proposalId, err)
	}

	return nil
}

func preMarketDeployProposals(environment string, currentNetworkParams *types.NetworkParams) ([]*commandspb.ProposalSubmission, error) {
	var commonProposals []*commandspb.ProposalSubmission

	isTradingEnabled, ok := currentNetworkParams.Params[netparams.PerpsMarketTradingEnabled]
	if !ok || isTradingEnabled != "1" {
		commonProposals = append(commonProposals, networkparameters.NewUpdateParameterProposalWithoutTime(netparams.PerpsMarketTradingEnabled, "1"))
	}

	currentL2Config, err := currentNetworkParams.GetEthereumL2Configs()
	if err != nil {
		return nil, fmt.Errorf("faled to get eth l2 config from network params: %w", err)
	}

	newL2Config := networkparameters.CloneEthereumL2Config(currentL2Config)

	for _, l2Config := range l2Configs[environment] {
		newL2Config, err = networkparameters.AppendEthereumL2Config(newL2Config, l2Config, true)
		if err != nil {
			return nil, fmt.Errorf("failed to append ethereum sepolia config to the l2 config: %w", err)
		}
	}

	if len(l2Configs[environment]) > 0 {
		l2ConfigJSON, err := json.Marshal(newL2Config)
		if err != nil {
			return nil, fmt.Errorf("failed to convert l2 config from proto to json: %w", err)
		}
		currentL2ConfigString, ok := currentNetworkParams.Params[netparams.BlockchainsEthereumL2Configs]
		if !ok || string(l2ConfigJSON) != currentL2ConfigString {
			commonProposals = append(commonProposals, networkparameters.NewUpdateParameterProposalWithoutTime(netparams.BlockchainsEthereumL2Configs, string(l2ConfigJSON)))
		}
	}

	return commonProposals, nil
}
