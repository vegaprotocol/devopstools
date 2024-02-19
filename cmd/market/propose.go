package market

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"code.vegaprotocol.io/vega/core/netparams"
	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/governance/market"
	"github.com/vegaprotocol/devopstools/governance/networkparameters"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"go.uber.org/zap"
)

const OracleAll = "*"

type ProposeArgs struct {
	*MarketArgs

	ProposeAAPL                        bool
	ProposeAAVEDAI                     bool
	ProposeBTCUSD                      bool
	ProposeETHBTC                      bool
	ProposeTSLA                        bool
	ProposeUNIDAI                      bool
	ProposeETHDAI                      bool
	ProposePerpetualBTCUSD             bool
	ProposePerpetualBTCUSDGnosis       bool
	ProposePerpetualSNXUSDUniswap      bool
	ProposePerpetualINJUSDUniswap      bool
	ProposePerpetual1000PEPEUSDUniswap bool
	ProposePerpetualLINKUSD            bool
	ProposePerpetualDAIUSD             bool
	ProposePerpetualEURUSD             bool
	ProposePerpetualETHUSD             bool
	ProposePerpetualJUPUSDPyth         bool
	ProposePerpetualBTCUSDPyth         bool

	ProposePerpetualBTCUSDMainnet  bool
	ProposePerpetualBTCUSDMainnet2 bool

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
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeETHDAI, "ethdai", false, "Propose ETHDAI market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposePerpetualBTCUSD, "perp-btcusd", false, "Propose perpetual BTCUSD market")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeAll, "all", false, "Propose all markets")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeCommunity, "community", false, "Propose community markets(only to devnet1)")

	proposeCmd.PersistentFlags().StringVar(&proposeArgs.OraclePubKey, "oracle-pubkey", "", "Oracle PubKey. Optional, by default proposer")
}

type MarketFlags struct {
	TotalMarkets int

	AAPL                        bool
	AAVEDAI                     bool
	BTCUSD                      bool
	ETHBTC                      bool
	TSLA                        bool
	UNIDAI                      bool
	ETHDAI                      bool
	CommunityLinkUSD            bool
	CommunityETHUSD             bool
	CommunityBTCUSD             bool
	PerpetualBTCUSD             bool
	PerpetualBTCUSDGnosis       bool
	PerpetualSNXUSDUniswap      bool
	PerpetualINJUSDUniswap      bool
	Perpetual1000PEPEUSDUniswap bool
	PerpetualDAIUSD             bool
	PerpetualEURUSD             bool
	PerpetualLINKUSD            bool
	PerpetualETHUSD             bool
	IncentiveBTCUSD             bool

	PerpetualJUPUSDPyth     bool
	PerpetualBTCUSDPyth     bool
	PerpetualBTCUSDMainnet  bool
	PerpetualBTCUSDMainnet2 bool

	MainnetBTCUSDT  bool
	MainnetDOGEUSDT bool
	MainnetETHUSDT  bool
	MainnetLINKUSDT bool
}

func dispatchMarkets(env string, args ProposeArgs) MarketFlags {
	// For fairground We want to keep clean network
	if env == types.NetworkFairground {
		return MarketFlags{
			TotalMarkets: 6,

			PerpetualBTCUSD: true,
			PerpetualETHUSD: true,

			MainnetBTCUSDT:              true,
			MainnetDOGEUSDT:             true,
			MainnetETHUSDT:              true,
			MainnetLINKUSDT:             true,
			PerpetualBTCUSDGnosis:       true,
			PerpetualSNXUSDUniswap:      true,
			PerpetualINJUSDUniswap:      true,
			Perpetual1000PEPEUSDUniswap: true,

			PerpetualJUPUSDPyth: true,
			PerpetualBTCUSDPyth: true,
		}
	}

	result := MarketFlags{
		// AAPL:                        args.ProposeAAPL || args.ProposeAll,
		// AAVEDAI:                     args.ProposeAAVEDAI || args.ProposeAll,
		// BTCUSD:                      args.ProposeBTCUSD || args.ProposeAll,
		// ETHBTC:                      args.ProposeETHBTC || args.ProposeAll,
		// TSLA:                        args.ProposeTSLA || args.ProposeAll,
		// UNIDAI:                      args.ProposeUNIDAI || args.ProposeAll,
		// ETHDAI:                      args.ProposeETHDAI || args.ProposeAll,
		PerpetualBTCUSDGnosis: args.ProposePerpetualBTCUSDGnosis || args.ProposeAll,
		// PerpetualSNXUSDUniswap:      args.ProposePerpetualSNXUSDUniswap || args.ProposeAll,
		// PerpetualINJUSDUniswap:      args.ProposePerpetualINJUSDUniswap || args.ProposeAll,
		// Perpetual1000PEPEUSDUniswap: args.ProposePerpetual1000PEPEUSDUniswap || args.ProposeAll,
		PerpetualBTCUSDMainnet:  args.ProposePerpetualBTCUSDMainnet || args.ProposeAll,
		PerpetualBTCUSDMainnet2: args.ProposePerpetualBTCUSDMainnet2 || args.ProposeAll,
	}

	if env == types.NetworkDevnet1 {
		result.CommunityBTCUSD = args.ProposeCommunity || args.ProposeAll
		result.CommunityETHUSD = args.ProposeCommunity || args.ProposeAll
		result.CommunityBTCUSD = args.ProposeCommunity || args.ProposeAll
		result.PerpetualBTCUSDGnosis = false  // We do not have L2 setup on devnet
		result.PerpetualSNXUSDUniswap = false // We do not have L2 setup on devnet
	}

	if env == types.NetworkDevnet1 || env == types.NetworkStagnet1 {
		result.PerpetualBTCUSD = args.ProposePerpetualBTCUSD || args.ProposeAll
		result.PerpetualEURUSD = args.ProposePerpetualEURUSD || args.ProposeAll
		result.PerpetualDAIUSD = args.ProposePerpetualDAIUSD || args.ProposeAll
		result.PerpetualETHUSD = args.ProposePerpetualETHUSD || args.ProposeAll
		result.PerpetualLINKUSD = args.ProposePerpetualLINKUSD || args.ProposeAll
	}

	result.TotalMarkets = tools.StructSize(result) - 1

	return result
}

func networkParametersForMarketPropose(network *veganetwork.VegaNetwork) (map[string]string, error) {
	statistics, err := network.DataNodeClient.Statistics()
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	networkVersion := semver.MustParse(statistics.Statistics.AppVersion)

	result := map[string]string{}

	prePerpsVersion := semver.New(0, 72, 99, "", "")

	if prePerpsVersion.Compare(networkVersion) <= 0 {
		result[netparams.PerpsMarketTradingEnabled] = "1"
	}

	preEthereumL2Version := semver.New(0, 73, 99, "", "")

	if preEthereumL2Version.Compare(networkVersion) <= 0 {
		if err := network.RefreshNetworkParams(); err != nil {
			return nil, fmt.Errorf("failed to refresh the network parameters: %w", err)
		}

		currentL2Config, err := network.NetworkParams.GetEthereumL2Configs()
		if err != nil {
			return nil, fmt.Errorf("failed to get the ethereum l2 configs: %w", err)
		}

		newL2Config := networkparameters.CloneEthereumL2Config(currentL2Config)
		l2ConfigChanged := false

		for _, l2Config := range l2Configs[network.Network] {
			l2ConfigChanged = true
			newL2Config, err = networkparameters.AppendEthereumL2Config(newL2Config, l2Config, true)
			if err != nil {
				return nil, fmt.Errorf("failed to append ethereum sepolia config to the l2 config: %w", err)
			}

		}

		if l2ConfigChanged {
			l2ConfigJSON, err := json.Marshal(newL2Config)
			if err != nil {
				return nil, fmt.Errorf("failed to convert l2 config from proto to json: %w", err)
			}
			result[netparams.BlockchainsEthereumL2Configs] = string(l2ConfigJSON)
		}
	}

	return result, nil
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

	assets, err := network.DataNodeClient.GetAssets()
	if err != nil {
		return fmt.Errorf("failed to get assets: %w", err)
	}

	settlementAssetId, foundSettlementAssetId := settlementAssetIDs[args.VegaNetworkName]
	if !foundSettlementAssetId {
		return fmt.Errorf("failed to get assets id's for network %s", err)
	}

	marketsFlags := dispatchMarkets(network.Network, args)

	// Propose
	resultsChannel := make(chan error, marketsFlags.TotalMarkets)
	var wg sync.WaitGroup

	networkParametersToUpdate, err := networkParametersForMarketPropose(network)
	if err != nil {
		return fmt.Errorf("failed to get network parameters that needs to be updated: %w", err)
	}
	if len(networkParametersToUpdate) > 0 {
		args.Logger.Sugar().Infof("Voting network parmeters required for markets creation: %v", networkParametersToUpdate)

		if _, err := governance.ProposeAndVoteOnNetworkParameters(
			networkParametersToUpdate, network.VegaTokenWhale, network.NetworkParams, network.DataNodeClient, args.Logger,
		); err != nil {
			return fmt.Errorf("failed to update network parameters required for market creation: %w", err)
		}

		time.Sleep(5 * time.Second)
		args.Logger.Info("Network parameters updated")
	}

	return nil

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
				closingTime, enactmentTime, MARKET_AAPL_MARKER, sub, logger,
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
				closingTime, enactmentTime, MARKET_AAVEDAI_MARKER, sub, logger,
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
				closingTime, enactmentTime, MARKET_BTCUSD_MARKER, sub, logger,
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
				closingTime, enactmentTime, MARKET_ETHBTC_MARKER, sub, logger,
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
				closingTime, enactmentTime, MARKET_TSLA_MARKER, sub, logger,
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
				closingTime, enactmentTime, MARKET_UNIDAI_MARKER, sub, logger,
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
				closingTime, enactmentTime, MARKET_ETHDAI_MARKER, sub, logger,
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
				closingTime, enactmentTime, market.CommunityBTCUSD230630MetadataID, sub, logger,
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
				closingTime, enactmentTime, market.CommunityETHUSD230630MetadataID, sub, logger,
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
				closingTime, enactmentTime, market.CommunityLinkUSD230630MetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualBTCUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewBTCUSDPerpetualMarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT,
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
				closingTime, enactmentTime, market.PerpetualBTCUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualBTCUSDGnosis {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewBTCUSDGnosisPerpetualMarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT,
				market.PerpetualBTCUSDOracleGnosisAddress,
				closingTime, enactmentTime,
				[]string{market.PerpetualBTCUSDGnosis},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual BTC USD (gnosis fake oracle)", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.PerpetualBTCUSDGnosis, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualSNXUSDUniswap {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewSNXUSDPerpetualMarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT,
				closingTime, enactmentTime,
				[]string{market.PerpetualSNXUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			// FIXME: I removed the need for the address because it's the hardcoded uniswap one???
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual SNX USD (mainnet uniswap oracle)", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				/* vvvv this vvv                       */
				closingTime, enactmentTime, market.PerpetualSNXUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualINJUSDUniswap {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewINJUSDPerpetualMarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT,
				closingTime, enactmentTime,
				[]string{market.PerpetualINJUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual INJ USD (mainnet uniswap oracle)", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.PerpetualINJUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.Perpetual1000PEPEUSDUniswap {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.New1000PEPEUSDPerpetualMarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT,
				closingTime, enactmentTime,
				[]string{market.Perpetual1000PEPEUSD},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual 1000PEPE USD (mainnet uniswap oracle)", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.Perpetual1000PEPEUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualJUPUSDPyth {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewJUPUSDPythPerpetualMarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT,
				closingTime, enactmentTime,
				[]string{market.PerpetualJUPUSDPyth},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			// FIXME: I removed the need for the address because it's the hardcoded uniswap one???
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual SNX USD (mainnet uniswap oracle)", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.PerpetualJUPUSDPyth, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualBTCUSDPyth {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewBTCUSDPythPerpetualMarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT,
				closingTime, enactmentTime,
				[]string{market.PerpetualBTCUSDPyth},
			)
			if err != nil {
				resultsChannel <- err
				return
			}
			// FIXME: I removed the need for the address because it's the hardcoded uniswap one???
			resultsChannel <- governance.ProposeVoteProvideLP(
				"Perpetual SNX USD (mainnet uniswap oracle)", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.PerpetualBTCUSDPyth, sub, logger,
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
				closingTime, enactmentTime, market.PerpetualLINKUSD, sub, logger,
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
				closingTime, enactmentTime, market.PerpetualDAIUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualETHUSD {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewETHUSDPerpetualMarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT,
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
				closingTime, enactmentTime, market.PerpetualETHUSD, sub, logger,
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
				closingTime, enactmentTime, market.PerpetualEURUSD, sub, logger,
			)
		}()
	}

	if marketsFlags.IncentiveBTCUSD {
		if _, assetExists := assets[market.IncentiveVegaAssetId]; !assetExists {
			logger.Warn(fmt.Sprintf("Cannot create incentive market. The %s asset does not exist on the network", market.IncentiveVegaAssetId))
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sub := market.NewBTCUSDIncentiveMarketProposal(
					market.IncentiveVegaAssetId, 5,
					market.IncentiveBTCUSDOracleAddress,
					closingTime, enactmentTime,
					[]string{market.IncentiveBTCUSD},
				)
				if err != nil {
					resultsChannel <- err
					return
				}
				resultsChannel <- governance.ProposeVoteProvideLP(
					"Incentive BTC USD", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
					closingTime, enactmentTime, market.IncentiveBTCUSD, sub, logger,
				)
			}()
		}
	}

	if marketsFlags.MainnetBTCUSDT {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewMainnetBTCUSDT(
				settlementAssetId.MainnetLikeAsset_USDT, CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
				"BTC USDT", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.MainnetBTCUSDTMetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.MainnetDOGEUSDT {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewMainnetDogeUSDT(
				settlementAssetId.MainnetLikeAsset_USDT, CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
				"DOGE USDT", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.MainnetDogeUSDTMetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.MainnetETHUSDT {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewMainnetETHUSDT(
				settlementAssetId.MainnetLikeAsset_USDT, CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
				"ETH USDT", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.MainnetETHUSDTMetadataID, sub, logger,
			)
		}()
	}

	if marketsFlags.MainnetLINKUSDT {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewMainnetLinkUSDT(
				settlementAssetId.MainnetLikeAsset_USDT, CoinBaseOraclePubKey,
				closingTime, enactmentTime,
				[]string{},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
				"LINK USDT", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.MainnetLinkUSDTMetadataID, sub, logger,
			)
		}()
	}

	// ---

	if marketsFlags.PerpetualBTCUSDMainnet {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewBTCUSDMainnetMarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT, closingTime, enactmentTime,
				[]string{},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
				"BTC USD Mainnet", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.EmptyMetadata, sub, logger,
			)
		}()
	}

	if marketsFlags.PerpetualBTCUSDMainnet2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := market.NewBTCUSDMainnet2MarketProposal(
				settlementAssetId.MainnetLikeAsset_USDT, closingTime, enactmentTime,
				[]string{},
			)
			resultsChannel <- governance.ProposeVoteProvideLP(
				"BTC USD Mainnet 2", network.DataNodeClient, lastBlockData, markets, proposerVegawallet,
				closingTime, enactmentTime, market.EmptyMetadata, sub, logger,
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
