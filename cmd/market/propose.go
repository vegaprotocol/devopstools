package market

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/governance/networkparameters"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"code.vegaprotocol.io/vega/core/netparams"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"
	"code.vegaprotocol.io/vega/wallet/wallet"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ProposeArgs struct {
	*Args
}

var proposeArgs ProposeArgs

var proposeCmd = &cobra.Command{
	Use:   "propose",
	Short: "Propose and vote on market",
	Long:  "Propose and vote on market",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPropose(proposeArgs); err != nil {
			proposeArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	proposeArgs.Args = &args

	Cmd.AddCommand(proposeCmd)
}

func runPropose(args ProposeArgs) error {
	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Info("Network file loaded", zap.String("name", cfg.Name.String()))

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Info("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Info("Looking for healthy gRPC endpoints...")
	healthyEndpoints := tools.FilterHealthyGRPCEndpoints(endpoints)
	if len(healthyEndpoints) == 0 {
		return fmt.Errorf("no healthy gRPC endpoint found on configured datanodes")
	}
	logger.Info("Healthy gRPC endpoints found", zap.Strings("endpoints", healthyEndpoints))

	datanodeClient := datanode.New(healthyEndpoints, 3*time.Second, args.Logger.Named("datanode"))

	logger.Info("Connecting to a datanode's gRPC endpoint...")
	dialCtx, cancelDialing := context.WithTimeout(ctx, 2*time.Second)
	defer cancelDialing()
	datanodeClient.MustDialConnection(dialCtx) // blocking
	logger.Info("Connected to a datanode's gRPC node", zap.String("node", datanodeClient.Target()))

	whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	if err := vega.GenerateKeysUpToKey(whaleWallet, cfg.Network.Wallets.VegaTokenWhale.PublicKey); err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	whalePublicKey := cfg.Network.Wallets.VegaTokenWhale.PublicKey

	logger.Info("Retrieving network parameters...")
	networkParams, err := datanodeClient.ListNetworkParameters(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}
	logger.Info("Network parameters retrieved")

	networkParamsProposals, err := preMarketDeployProposals(cfg.Name, networkParams)
	if err != nil {
		return fmt.Errorf("failed to determine network parameters to update before markets are proposed: %w", err)
	}

	if len(networkParamsProposals) > 0 {
		logger.Info("Updating network parameters...")
		networkParamsBatchProposal := governance.NewBatchProposal(
			fmt.Sprintf("%q devopstools network params proposal", cfg.Name.String()),
			"Update network parameters before markets are proposed",
			time.Now().Add(30*time.Second),
			networkParamsProposals,
			nil,
		)

		if err := sendBatchProposal(ctx, logger, datanodeClient, whaleWallet, whalePublicKey, networkParamsBatchProposal); err != nil {
			return fmt.Errorf("failed to send batch proposal to update network parameters: %w", err)
		}
		logger.Info("Network parameters updated")
	} else {
		logger.Info("Network parameters do not need to be updated")
	}

	allMarkets, err := datanodeClient.GetAllMarketsWithState(ctx, datanode.ActiveMarkets)
	if err != nil {
		return fmt.Errorf("could not retrieve markets from datanode: %w", err)
	}

	missingMarketsProposals := collectMissingMarkets(allMarkets, ProposalsForEnvironment(cfg.Name), logger)

	if len(missingMarketsProposals) < 1 {
		logger.Info("No market to propose")
		return nil
	}

	marketsBatchProposal := governance.NewBatchProposal(
		fmt.Sprintf("%q devopstools markets proposal", cfg.Name.String()),
		"Create all markets managed by devopstools",
		time.Now().Add(30*time.Second),
		missingMarketsProposals,
		nil,
	)

	if len(missingMarketsProposals) < 1 {
		logger.Info("All required markets exist")
		return nil
	}

	return sendBatchProposal(ctx, logger, datanodeClient, whaleWallet, whalePublicKey, marketsBatchProposal)
}

func collectMissingMarkets(allMarkets []*vegapb.Market, allNetworkProposals []*commandspb.ProposalSubmission, logger *zap.Logger) []*commandspb.ProposalSubmission {
	var result []*commandspb.ProposalSubmission

	for idx, proposal := range allNetworkProposals {
		newMarketProposal, ok := proposal.Terms.Change.(*vegapb.ProposalTerms_NewMarket)
		if !ok {
			continue
		}

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
			logger.Info("Adding proposal for market creation to batch", zap.String("market-code", marketCode))
		} else {
			logger.Info("Market already existing on the network", zap.String("market-code", marketCode))
		}
	}

	return result
}

func sendBatchProposal(ctx context.Context, logger *zap.Logger, datanodeClient *datanode.DataNode, whaleWallet wallet.Wallet, whaleKey string, proposals *commandspb.BatchProposalSubmission) error {
	walletTxReq := walletpb.SubmitTransactionRequest{
		Command: &walletpb.SubmitTransactionRequest_BatchProposalSubmission{
			BatchProposalSubmission: proposals,
		},
	}

	logger.Info("Sending transaction to the network")
	resp, err := walletpkg.SendTransaction(ctx, whaleWallet, whaleKey, &walletTxReq, datanodeClient)
	if err != nil {
		return fmt.Errorf("failed to submit batch proposal with signature: %w", err)
	}
	logger.Sugar().Infof("Batch proposal transaction ID: %s", resp.TxHash)

	logger.Info("Searching proposal ID")
	proposalId, err := tools.RetryReturn(5, time.Second*5, func() (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		return governance.FindProposalID(ctx, whaleKey, proposals.Reference, datanodeClient)
	})
	if err != nil {
		return fmt.Errorf("failed to find proposal ID: %w", err)
	}
	logger.Sugar().Info("Found proposal with the following ID: %s", proposalId)

	if err = tools.RetryRun(10, 6*time.Second, func() error {
		return governance.VoteOnProposal(ctx, "BatchProposal vote", proposalId, whaleWallet, whaleKey, datanodeClient)
	}); err != nil {
		return fmt.Errorf("failed to vote for batch proposal %q: %w", proposalId, err)
	}

	return nil
}

func preMarketDeployProposals(environment config.NetworkName, currentNetworkParams *types.NetworkParams) ([]*commandspb.ProposalSubmission, error) {
	var commonProposals []*commandspb.ProposalSubmission

	isTradingEnabled, ok := currentNetworkParams.Params[netparams.PerpsMarketTradingEnabled]
	if !ok || isTradingEnabled != "1" {
		commonProposals = append(commonProposals, networkparameters.NewUpdateParameterProposalWithoutTime(netparams.PerpsMarketTradingEnabled, "1"))
	}

	currentL2Config, err := currentNetworkParams.GetEthereumL2Configs()
	if err != nil {
		return nil, fmt.Errorf("could not get Ethereum L2 configuration from network paramters: %w", err)
	}

	newL2Config := networkparameters.CloneEthereumL2Config(currentL2Config)

	for _, l2Config := range l2Configs[environment] {
		newL2Config, err = networkparameters.AppendEthereumL2Config(newL2Config, l2Config, true)
		if err != nil {
			return nil, fmt.Errorf("could not add new Ethereum L2 configuration to existing ones: %w", err)
		}
	}

	if len(l2Configs[environment]) > 0 {
		l2ConfigJSON, err := json.Marshal(newL2Config)
		if err != nil {
			return nil, fmt.Errorf("could not serialize Ethereum L2 configuration to JSON: %w", err)
		}
		currentL2ConfigString, ok := currentNetworkParams.Params[netparams.BlockchainsEthereumL2Configs]
		if !ok || string(l2ConfigJSON) != currentL2ConfigString {
			commonProposals = append(commonProposals, networkparameters.NewUpdateParameterProposalWithoutTime(netparams.BlockchainsEthereumL2Configs, string(l2ConfigJSON)))
		}
	}

	return commonProposals, nil
}
