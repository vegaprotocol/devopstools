package market

import (
	"fmt"
	"os"
	"sync"
	"time"

	"code.vegaprotocol.io/vega/core/netparams"
	vegatypes "code.vegaprotocol.io/vega/core/types"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/proposals"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type ProposeArgs struct {
	*MarketArgs
	ProposeAAPL    bool
	ProposeAAVEDAI bool
	ProposeBTCUSD  bool
	ProposeETHBTC  bool
	ProposeTSLA    bool
	ProposeUNIDAI  bool
	ProposeETHDAI  bool
	ProposeAll     bool
	OraclePubKey   string
	FakeAsset      bool
	ERC20Asset     bool
}

var proposeArgs ProposeArgs

// proposeCmd represents the propose command
var proposeCmd = &cobra.Command{
	Use:   "propose",
	Short: "Propose and vote on market",
	Long:  `Propose and vote on market`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunPropose(proposeArgs); err != nil {
			proposeArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	proposeArgs.MarketArgs = &marketArgs

	MarketCmd.AddCommand(proposeCmd)

	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeAAPL, "aapl", false, "Propose AAPL market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeAAVEDAI, "aavedai", false, "Propose AAVEDAI market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeBTCUSD, "btcusd", false, "Propose BTCUSD market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeETHBTC, "ethbtc", false, "Propose ETHBTC market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeTSLA, "tsla", false, "Propose TSLA market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeUNIDAI, "unidai", false, "Propose UNIDAI market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeETHDAI, "ethdai", false, "Propose ETHDI market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeAll, "all", false, "Propose all markets")

	proposeCmd.PersistentFlags().StringVar(&proposeArgs.OraclePubKey, "oracle-pubkey", "", "Oracle PubKey. Optional, by default proposer")
}

func RunPropose(args ProposeArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()
	minClose, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinClose])
	if err != nil {
		return err
	}
	minEnact, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinEnact])
	if err != nil {
		return err
	}

	var (
		proposerVegawallet = network.VegaTokenWhale
		oraclePubKey       = args.OraclePubKey
		closingTime        = time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime      = time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)
		logger             = args.Logger
	)
	if len(oraclePubKey) == 0 {
		oraclePubKey = proposerVegawallet.PublicKey
	}
	lastBlockData, err := network.DataNodeClient.LastBlockData()
	if err != nil {
		return err
	}
	markets, err := network.DataNodeClient.GetAllMarkets()
	if err != nil {
		return err
	}

	settlementAssetId, foundSettlementAssetId := settlementAssetIDs[args.VegaNetworkName]
	if !foundSettlementAssetId {
		return fmt.Errorf("failed to get assets id's for network %s", err)
	}

	// Propose
	resultsChannel := make(chan error, 10)
	var wg sync.WaitGroup

	//
	// AAPL
	//
	if args.ProposeAAPL || args.ProposeAll {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewAAPLMarketProposal(
				settlementAssetId.AAPL, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_AAPL_MARKER},
			)
			resultsChannel <- proposeVoteProvideLP(
				"AAPL", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				oraclePubKey, closingTime, enactmentTime, MARKET_AAPL_MARKER, sub, logger,
			)
		}()
	}
	//
	// AAVEDAI
	//
	if args.ProposeAAVEDAI || args.ProposeAll {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewAAVEDAIMarketProposal(
				settlementAssetId.AAVEDAI, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_AAVEDAI_MARKER},
			)
			resultsChannel <- proposeVoteProvideLP(
				"AAVEDAI", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				oraclePubKey, closingTime, enactmentTime, MARKET_AAVEDAI_MARKER, sub, logger,
			)
		}()
	}
	//
	// BTCUSD
	//
	if args.ProposeBTCUSD || args.ProposeAll {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewBTCUSDMarketProposal(
				settlementAssetId.BTCUSD, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_BTCUSD_MARKER},
			)
			resultsChannel <- proposeVoteProvideLP(
				"BTCUSD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				oraclePubKey, closingTime, enactmentTime, MARKET_BTCUSD_MARKER, sub, logger,
			)
		}()
	}

	//
	// ETHBTC
	//
	if args.ProposeETHBTC || args.ProposeAll {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewETHBTCMarketProposal(
				settlementAssetId.ETHBTC, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_ETHBTC_MARKER},
			)
			resultsChannel <- proposeVoteProvideLP(
				"ETHBTC", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				oraclePubKey, closingTime, enactmentTime, MARKET_ETHBTC_MARKER, sub, logger,
			)
		}()
	}
	//
	// TSLA
	//
	if args.ProposeTSLA || args.ProposeAll {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewTSLAMarketProposal(
				settlementAssetId.TSLA, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_TSLA_MARKER},
			)
			resultsChannel <- proposeVoteProvideLP(
				"TSLA", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				oraclePubKey, closingTime, enactmentTime, MARKET_TSLA_MARKER, sub, logger,
			)
		}()
	}
	//
	// UNIDAI
	//
	if args.ProposeUNIDAI || args.ProposeAll {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewUNIDAIMarketProposal(
				settlementAssetId.UNIDAI, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_UNIDAI_MARKER},
			)
			resultsChannel <- proposeVoteProvideLP(
				"UNIDAI", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				oraclePubKey, closingTime, enactmentTime, MARKET_UNIDAI_MARKER, sub, logger,
			)
		}()
	}

	//
	// ETHDAI
	//
	if args.ProposeETHDAI || args.ProposeAll {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewETHDAIMarketProposal(
				settlementAssetId.ETHDAI, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_ETHDAI_MARKER},
			)
			resultsChannel <- proposeVoteProvideLP(
				"ETHDAI", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				oraclePubKey, closingTime, enactmentTime, MARKET_ETHDAI_MARKER, sub, logger,
			)
		}()
	}

	wg.Wait()
	close(resultsChannel)

	for err := range resultsChannel {
		if err != nil {
			return fmt.Errorf("at least one proposal failed: %w", err)
		}
	}

	return nil
}

func getMarket(markets []*vega.Market, oraclePubKey string, metadataTag string) *vega.Market {
	for _, market := range markets {
		if market.TradableInstrument == nil || market.TradableInstrument.Instrument == nil ||
			market.TradableInstrument.Instrument.GetFuture() == nil ||
			market.TradableInstrument.Instrument.GetFuture().DataSourceSpecForTradingTermination == nil ||
			market.TradableInstrument.Instrument.GetFuture().DataSourceSpecForTradingTermination.Config == nil {
			continue
		}

		signers := market.TradableInstrument.Instrument.GetFuture().DataSourceSpecForTradingTermination.Config.Signers
		stringSigners := []string{}
		for _, signer := range signers {
			if signer.GetPubKey() == nil {
				continue
			}

			stringSigners = append(stringSigners, signer.GetPubKey().GetKey())
		}

		if slices.Contains(
			stringSigners,
			oraclePubKey,
		) && slices.Contains(
			market.TradableInstrument.Instrument.Metadata.Tags,
			metadataTag,
		) {
			// TODO check if open
			return market
		}
	}
	return nil
}

func proposeVoteProvideLP(
	name string,
	dataNodeClient vegaapi.DataNodeClient,
	lastBlockData *vegaapipb.LastBlockHeightResponse,
	markets []*vega.Market,
	proposerVegawallet *wallet.VegaWallet,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	marketMetadataMarker string,
	proposal *commandspb.ProposalSubmission,
	logger *zap.Logger,
) error {
	market := getMarket(markets, oraclePubKey, marketMetadataMarker)
	if market != nil {
		logger.Info("market already exist", zap.String("market", market.Id))
		return nil
	}
	reference := proposal.Reference
	//
	// PROPOSE
	//
	// Prepare vegawallet Transaction Request
	walletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_ProposalSubmission{
			ProposalSubmission: proposal,
		},
	}
	if err := submitTx("propose market", dataNodeClient, proposerVegawallet, logger, &walletTxReq); err != nil {
		return err
	}

	//
	// Find Proposal
	//
	time.Sleep(time.Second * 5)
	res, err := dataNodeClient.ListGovernanceData(&v2.ListGovernanceDataRequest{
		ProposalReference: &reference,
	})
	if err != nil {
		return err
	}
	var proposalId string
	for _, edge := range res.Connection.Edges {
		if edge.Node.Proposal.Reference == reference {
			logger.Info("Found proposal", zap.String("market", name), zap.String("reference", reference),
				zap.String("status", edge.Node.Proposal.State.String()),
				zap.Any("proposal", edge.Node.Proposal))
			proposalId = edge.Node.Proposal.Id
		}
	}
	//
	// VOTE
	//
	voteWalletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
			VoteSubmission: &commandspb.VoteSubmission{
				ProposalId: proposalId,
				Value:      vega.Vote_VALUE_YES,
			},
		},
	}
	if err := submitTx("vote on market proposal", dataNodeClient, proposerVegawallet, logger, &voteWalletTxReq); err != nil {
		return err
	}

	//
	// Provide LP
	//

	markets, err = dataNodeClient.GetAllMarkets()
	if err != nil {
		return fmt.Errorf("failed to get markets")
	}
	market = getMarket(markets, oraclePubKey, marketMetadataMarker)
	if market == nil {
		return fmt.Errorf("failed to find particular market for %s oracle public key", oraclePubKey)
	}

	// Prepare vegawallet Transaction Request
	provideLPWalletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_LiquidityProvisionSubmission{
			LiquidityProvisionSubmission: &commandspb.LiquidityProvisionSubmission{
				Fee:              "0.01",
				MarketId:         market.Id,
				CommitmentAmount: "1",
				Buys: []*vega.LiquidityOrder{
					{
						Reference:  vegatypes.PeggedReferenceBestBid,
						Proportion: 10,
						Offset:     "1000",
					},
					{
						Reference:  vegatypes.PeggedReferenceBestBid,
						Proportion: 10,
						Offset:     "2000",
					},
				},
				Sells: []*vega.LiquidityOrder{
					{
						Reference:  vegatypes.PeggedReferenceBestAsk,
						Proportion: 10,
						Offset:     "2000",
					},
					{
						Reference:  vegatypes.PeggedReferenceBestAsk,
						Proportion: 10,
						Offset:     "1000",
					},
				},
			},
		},
	}
	if err := submitTx("Provide LP", dataNodeClient, proposerVegawallet, logger, &provideLPWalletTxReq); err != nil {
		return err
	}

	return nil
}

func submitTx(
	description string,
	dataNodeClient vegaapi.DataNodeClient,
	proposerVegawallet *wallet.VegaWallet,
	logger *zap.Logger,
	walletTxReq *walletpb.SubmitTransactionRequest,
) error {
	lastBlockData, err := dataNodeClient.LastBlockData()
	if err != nil {
		return fmt.Errorf("failed to submit tx, cos")
	}

	// Sign + Proof of Work vegawallet Transaction request
	signedTx, err := proposerVegawallet.SignTxWithPoW(walletTxReq, lastBlockData)
	if err != nil {
		logger.Error("Failed to sign a trasnaction", zap.String("description", description),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", &walletTxReq), zap.Error(err))
		return err
	}

	// wrap in vega Transaction Request
	submitReq := &vegaapipb.SubmitTransactionRequest{
		Tx:   signedTx,
		Type: vegaapipb.SubmitTransactionRequest_TYPE_SYNC,
	}

	// Submit Transaction
	logger.Info("Submit transaction", zap.String("description", description),
		zap.String("proposer", proposerVegawallet.PublicKey))
	submitResponse, err := dataNodeClient.SubmitTransaction(submitReq)
	if err != nil {
		logger.Error("Failed to submit a trasnaction", zap.String("description", description),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", submitReq), zap.Error(err))
		return err
	}
	if !submitResponse.Success {
		logger.Error("Transaction submission response is not successful",
			zap.String("proposer", proposerVegawallet.PublicKey), zap.String("description", description),
			zap.Any("txReq", submitReq), zap.Any("response", submitResponse))
		return err
	}
	logger.Info("Successful Submision of Market Proposal", zap.String("description", description),
		zap.String("proposer", proposerVegawallet.PublicKey), zap.String("txHash", submitResponse.TxHash),
		zap.Any("response", submitResponse))

	return nil
}
