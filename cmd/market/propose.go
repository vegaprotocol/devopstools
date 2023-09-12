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
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

const OracleAll = "*"

type ProposeArgs struct {
	*MarketArgs
	ProposeAAPL             bool
	ProposeAAVEDAI          bool
	ProposeBTCUSD           bool
	ProposeETHBTC           bool
	ProposeTSLA             bool
	ProposeUNIDAI           bool
	ProposeETHDAI           bool
	ProposePerpetualBTCUSD  bool
	ProposePerpetualLINKUSD bool
	ProposePerpetualDAIUSD  bool
	ProposePerpetualEURUSD  bool
	ProposePerpetualETHUSD  bool

	ProposeAll       bool
	ProposeCommunity bool

	OraclePubKey string
	FakeAsset    bool
	ERC20Asset   bool
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
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposePerpetualBTCUSD, "perp-btcusd", false, "Propose perpetual BTCUSD market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeAll, "all", false, "Propose all markets")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeCommunity, "community", false, "Propose community markets(only to devnet1)")

	proposeCmd.PersistentFlags().StringVar(&proposeArgs.OraclePubKey, "oracle-pubkey", "", "Oracle PubKey. Optional, by default proposer")
}

type MarketFlags struct {
	TotalMarkets int

	AAPL             bool
	AAVEDAI          bool
	BTCUSD           bool
	ETHBTC           bool
	TSLA             bool
	UNIDAI           bool
	ETHDAI           bool
	CommunityLinkUSD bool
	CommunityETHUSD  bool
	CommunityBTCUSD  bool
	PerpetualBTCUSD  bool
	PerpetualDAIUSD  bool
	PerpetualEURUSD  bool
	PerpetualLINKUSD bool
	PerpetualETHUSD  bool
}

func dispatchMarkets(env string, args ProposeArgs) MarketFlags {
	result := MarketFlags{
		AAPL:    args.ProposeAAPL || args.ProposeAll,
		AAVEDAI: args.ProposeAAVEDAI || args.ProposeAll,
		BTCUSD:  args.ProposeBTCUSD || args.ProposeAll,
		ETHBTC:  args.ProposeETHBTC || args.ProposeAll,
		TSLA:    args.ProposeTSLA || args.ProposeAll,
		UNIDAI:  args.ProposeUNIDAI || args.ProposeAll,
		ETHDAI:  args.ProposeETHDAI || args.ProposeAll,
	}

	if env == types.NetworkDevnet1 {
		result.CommunityBTCUSD = args.ProposeCommunity || args.ProposeAll
		result.CommunityETHUSD = args.ProposeCommunity || args.ProposeAll
		result.CommunityBTCUSD = args.ProposeCommunity || args.ProposeAll
	}

	if env == types.NetworkDevnet1 || env == types.Stagnet1 {
		result.PerpetualBTCUSD = args.ProposePerpetualBTCUSD || args.ProposeAll
		result.PerpetualEURUSD = args.ProposePerpetualEURUSD || args.ProposeAll
		result.PerpetualDAIUSD = args.ProposePerpetualDAIUSD || args.ProposeAll
		result.PerpetualETHUSD = args.ProposePerpetualETHUSD || args.ProposeAll
		result.PerpetualLINKUSD = args.ProposePerpetualLINKUSD || args.ProposeAll
	}

	result.TotalMarkets = tools.StructSize(result) - 1

	return result
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

	marketsFlags := dispatchMarkets(network.Network, args)

	// Propose
	resultsChannel := make(chan error, marketsFlags.TotalMarkets)
	var wg sync.WaitGroup

	//
	// AAPL
	//
	if marketsFlags.AAPL {
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
	if marketsFlags.AAVEDAI {
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
	if marketsFlags.BTCUSD {
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
	if marketsFlags.ETHBTC {
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
	if marketsFlags.TSLA {
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
	if marketsFlags.UNIDAI {
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
	if marketsFlags.ETHDAI {
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

	if marketsFlags.CommunityBTCUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, err := proposals.NewCommunityBTCUSD230630(
				settlementAssetId.SettlementAsset_USDC, 6,
				CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{proposals.CommunityBTCUSD230630MetadataID},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- proposeVoteProvideLP(
				"Community BTC USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				CoinBaseOraclePubKey, closingTime, enactmentTime, proposals.CommunityBTCUSD230630MetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.CommunityETHUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, err := proposals.NewCommunityETHUSD230630(
				settlementAssetId.SettlementAsset_USDC, 6,
				CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{proposals.CommunityETHUSD230630MetadataID},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- proposeVoteProvideLP(
				"Community ETH USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				CoinBaseOraclePubKey, closingTime, enactmentTime, proposals.CommunityETHUSD230630MetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.CommunityLinkUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, err := proposals.NewCommunityLinkUSD230630(
				settlementAssetId.SettlementAsset_USDC, 6,
				CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{proposals.CommunityLinkUSD230630MetadataID},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- proposeVoteProvideLP(
				"Community LINK USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				CoinBaseOraclePubKey, closingTime, enactmentTime, proposals.CommunityLinkUSD230630MetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualBTCUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewBTCUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 6,
				proposals.PerpetualBTCUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{proposals.PerpetualBTCUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- proposeVoteProvideLP(
				"Perpetual BTC USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				proposals.PerpetualBTCUSDOracleAddress, closingTime, enactmentTime, proposals.PerpetualBTCUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualLINKUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewLINKUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 6,
				proposals.PerpetualLINKUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{proposals.PerpetualLINKUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- proposeVoteProvideLP(
				"Perpetual LINK USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				proposals.PerpetualLINKUSDOracleAddress, closingTime, enactmentTime, proposals.PerpetualLINKUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualDAIUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewDAIUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 6,
				proposals.PerpetualDAIUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{proposals.PerpetualDAIUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- proposeVoteProvideLP(
				"Perpetual DAI USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				proposals.PerpetualDAIUSDOracleAddress, closingTime, enactmentTime, proposals.PerpetualDAIUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualETHUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewETHUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 6,
				proposals.PerpetualETHUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{proposals.PerpetualETHUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- proposeVoteProvideLP(
				"Perpetual ETH USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				proposals.PerpetualETHUSDOracleAddress, closingTime, enactmentTime, proposals.PerpetualETHUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualEURUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := proposals.NewEURUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 6,
				proposals.PerpetualEURUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{proposals.PerpetualEURUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- proposeVoteProvideLP(
				"Perpetual EUR USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				proposals.PerpetualEURUSDOracleAddress, closingTime, enactmentTime, proposals.PerpetualEURUSD, sub, logger,
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
		if market.TradableInstrument == nil || market.TradableInstrument.Instrument == nil {
			continue
		}

		future := market.TradableInstrument.Instrument.GetFuture()
		perpetual := market.TradableInstrument.Instrument.GetPerpetual()

		if future == nil && perpetual == nil {
			// TODO: Log the unsupported market here?
			continue
		}

		stringSigners := []string{OracleAll}
		if perpetual != nil &&
			perpetual.DataSourceSpecForSettlementData != nil &&
			perpetual.DataSourceSpecForSettlementData.GetData() != nil &&
			perpetual.DataSourceSpecForSettlementData.GetData().GetExternal() != nil {
			external := perpetual.DataSourceSpecForSettlementData.GetData().GetExternal()

			if external.GetOracle() != nil {
				for _, signer := range external.GetOracle().Signers {
					if signer.GetPubKey() != nil {
						stringSigners = append(stringSigners, signer.GetPubKey().GetKey())
					}
					if signer.GetEthAddress() != nil {
						stringSigners = append(stringSigners, signer.GetEthAddress().Address)
					}
				}
			}

			if external.GetEthOracle() != nil {
				stringSigners = append(stringSigners, external.GetEthOracle().Address)
			}
		}

		if future != nil &&
			future.DataSourceSpecForSettlementData != nil &&
			future.DataSourceSpecForSettlementData.GetData() != nil &&
			future.DataSourceSpecForSettlementData.GetData().GetExternal() != nil &&
			future.DataSourceSpecForSettlementData.GetData().GetExternal().GetOracle() != nil {
			signers := future.DataSourceSpecForSettlementData.GetData().GetExternal().GetOracle().Signers
			for _, signer := range signers {
				if signer.GetPubKey() != nil {
					stringSigners = append(stringSigners, signer.GetPubKey().GetKey())
				}
				if signer.GetEthAddress() != nil {
					stringSigners = append(stringSigners, signer.GetEthAddress().Address)
				}
			}
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
		logger.Info("market already exist", zap.String("market", market.Id), zap.String("name", name))
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
	time.Sleep(time.Second * 10)
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

	if len(proposalId) < 1 {
		return fmt.Errorf("got empty proposal id for the %s reference", reference)
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
		return fmt.Errorf("failed to get markets: %w", err)
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
		return fmt.Errorf("failed to submit tx: %w", err)
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
			zap.Any("txReq", submitReq.String()), zap.String("response", fmt.Sprintf("%#v", submitResponse)))
		return err
	}
	logger.Info("Successful Submision of Market Proposal", zap.String("description", description),
		zap.String("proposer", proposerVegawallet.PublicKey), zap.String("txHash", submitResponse.TxHash),
		zap.Any("response", submitResponse))

	return nil
}
