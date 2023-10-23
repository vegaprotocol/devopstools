package market

import (
	"fmt"
	"os"
	"sync"
	"time"

	"code.vegaprotocol.io/vega/core/netparams"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/governance/market"
	"github.com/vegaprotocol/devopstools/governance/networkparameters"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"go.uber.org/zap"
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

	if env == types.NetworkDevnet1 || env == types.NetworkStagnet1 || env == types.NetworkFairground {
		result.PerpetualBTCUSD = args.ProposePerpetualBTCUSD || args.ProposeAll
		result.PerpetualEURUSD = args.ProposePerpetualEURUSD || args.ProposeAll
		result.PerpetualDAIUSD = args.ProposePerpetualDAIUSD || args.ProposeAll
		result.PerpetualETHUSD = args.ProposePerpetualETHUSD || args.ProposeAll
		result.PerpetualLINKUSD = args.ProposePerpetualLINKUSD || args.ProposeAll
	}

	result.TotalMarkets = tools.StructSize(result) - 1

	return result
}

func updateNetworkParameters(closingTime, enactmentTime time.Time, networkParams map[string]string, networkVersion string) []*commandspb.ProposalSubmission {
	result := []*commandspb.ProposalSubmission{}

	perpetualEnabled, perpetualEnabledParamExist := networkParams[netparams.PerpsMarketTradingEnabled]
	prePerpsVersion := semver.New(0, 72, 99, "", "")

	if prePerpsVersion.Compare(semver.MustParse(networkVersion)) <= 0 && (!perpetualEnabledParamExist || perpetualEnabled != "1") {
		result = append(result, networkparameters.NewUpdateParametersProposal(netparams.PerpsMarketTradingEnabled, "1", closingTime, enactmentTime))
	}

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
	statistics, err := network.DataNodeClient.Statistics()
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

	networkParametersProposals := updateNetworkParameters(closingTime, enactmentTime, network.NetworkParams.Params, statistics.Statistics.AppVersion)

	for _, p := range networkParametersProposals {
		closingTime = time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime = time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)

		if err := governance.ProposeAndVote(logger, proposerVegawallet, network.DataNodeClient, p); err != nil {
			return fmt.Errorf("failed to submit network parameter change proposal: %w", err)
		}
	}

	// Check if all network params has been updated correctly
	if len(networkParametersProposals) > 0 {
		err := tools.RetryRun(6, 10*time.Second, func() error {
			if err := network.RefreshNetworkParams(); err != nil {
				return fmt.Errorf("failed to refresh network parameters: %w", err)
			}

			for _, proposal := range networkParametersProposals {
				changes := proposal.Terms.GetUpdateNetworkParameter()
				if changes == nil {
					return fmt.Errorf("invalid proposal for %s", proposal.Rationale.Title)
				}

				value, ok := network.NetworkParams.Params[changes.Changes.Key]
				if !ok {
					return fmt.Errorf("the %s network parameter not found yet on the network", changes.Changes.Key)
				}

				if value != changes.Changes.Value {
					return fmt.Errorf(
						"the %s network parameter not updated yet: expected value %s, current value: %s",
						changes.Changes.Key,
						changes.Changes.Value,
						value,
					)
				}
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("network parameters not updated yet: %w", err)
		}
	}

	closingTime = time.Now().Add(time.Second * 20).Add(minClose)
	enactmentTime = time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)
	//
	// AAPL
	//
	if marketsFlags.AAPL {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewAAPLMarketProposal(
				settlementAssetId.AAPL, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_AAPL_MARKER},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
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
			sub := market.NewAAVEDAIMarketProposal(
				settlementAssetId.AAVEDAI, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_AAVEDAI_MARKER},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
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
			sub := market.NewBTCUSDMarketProposal(
				settlementAssetId.BTCUSD, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_BTCUSD_MARKER},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
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
			sub := market.NewETHBTCMarketProposal(
				settlementAssetId.ETHBTC, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_ETHBTC_MARKER},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
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
			sub := market.NewTSLAMarketProposal(
				settlementAssetId.TSLA, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_TSLA_MARKER},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
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
			sub := market.NewUNIDAIMarketProposal(
				settlementAssetId.UNIDAI, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_UNIDAI_MARKER},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
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
			sub := market.NewETHDAIMarketProposal(
				settlementAssetId.ETHDAI, 5, oraclePubKey,
				closingTime, enactmentTime,
				[]string{MARKET_ETHDAI_MARKER},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
				"ETHDAI", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				oraclePubKey, closingTime, enactmentTime, MARKET_ETHDAI_MARKER, sub, logger,
			)
		}()
	}

	if marketsFlags.CommunityBTCUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, err := market.NewCommunityBTCUSD230630(
				settlementAssetId.SettlementAsset_USDC, 6,
				CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{market.CommunityBTCUSD230630MetadataID},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Community BTC USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				CoinBaseOraclePubKey, closingTime, enactmentTime, market.CommunityBTCUSD230630MetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.CommunityETHUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, err := market.NewCommunityETHUSD230630(
				settlementAssetId.SettlementAsset_USDC, 6,
				CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{market.CommunityETHUSD230630MetadataID},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Community ETH USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				CoinBaseOraclePubKey, closingTime, enactmentTime, market.CommunityETHUSD230630MetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.CommunityLinkUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, err := market.NewCommunityLinkUSD230630(
				settlementAssetId.SettlementAsset_USDC, 6,
				CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{market.CommunityLinkUSD230630MetadataID},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Community LINK USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				CoinBaseOraclePubKey, closingTime, enactmentTime, market.CommunityLinkUSD230630MetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualBTCUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewBTCUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 5,
				market.PerpetualBTCUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{market.PerpetualBTCUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual BTC USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				market.PerpetualBTCUSDOracleAddress, closingTime, enactmentTime, market.PerpetualBTCUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualLINKUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewLINKUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 5,
				market.PerpetualLINKUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{market.PerpetualLINKUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual LINK USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				market.PerpetualLINKUSDOracleAddress, closingTime, enactmentTime, market.PerpetualLINKUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualDAIUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewDAIUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 5,
				market.PerpetualDAIUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{market.PerpetualDAIUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual DAI USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				market.PerpetualDAIUSDOracleAddress, closingTime, enactmentTime, market.PerpetualDAIUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualETHUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewETHUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 5,
				market.PerpetualETHUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{market.PerpetualETHUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual ETH USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				market.PerpetualETHUSDOracleAddress, closingTime, enactmentTime, market.PerpetualETHUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualEURUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewEURUSDPerpetualMarketProposal(
				settlementAssetId.SettlementAsset_USDC, 5,
				market.PerpetualEURUSDOracleAddress,
				closingTime, enactmentTime,
				[]string{market.PerpetualEURUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual EUR USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				market.PerpetualEURUSDOracleAddress, closingTime, enactmentTime, market.PerpetualEURUSD, sub, logger,
			)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		sub := market.NewBTCUSDPerpetualMarketProposal_New(
			"dd20590509d30d20bdbbe64dc1090c1140c7690121a9b9940bc66f62dfa2e599", 5,
			market.PerpetualBTCUSDOracleAddress,
			closingTime, enactmentTime,
			[]string{market.PerpetualBTCUSDNew},
		)
		if err != nil {
			resultsChannel <- err
			return
		}
		resultsChannel <- governance.ProposeVoteProvideLP(
			"BTC USD-R Incentive I", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
			market.PerpetualBTCUSDOracleAddress, closingTime, enactmentTime, market.PerpetualBTCUSDNew, sub, logger,
		)
	}()

	wg.Wait()
	close(resultsChannel)

	for err := range resultsChannel {
		if err != nil {
			return fmt.Errorf("at least one proposal failed: %w", err)
		}
	}

	return nil
}
