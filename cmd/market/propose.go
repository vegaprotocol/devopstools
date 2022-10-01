package market

import (
	"fmt"
	"os"
	"sync"
	"time"

	"code.vegaprotocol.io/vega/core/netparams"
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
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ProposeAll, "all", false, "Propose all markets")

	proposeCmd.PersistentFlags().StringVar(&proposeArgs.OraclePubKey, "oracle-pubkey", "", "Oracle PubKey. Optional, by default proposer")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.FakeAsset, "fake-asset", false, "Use Fake Assets as Settlement Assets. By default use ERC20 tokens")
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.ERC20Asset, "erc20-asset", false, "Use ERC20 Tokens as Settlement Assets.")
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
	// Assets
	assets, err := network.DataNodeClient.GetAssets()
	if err != nil {
		return err
	}

	// Propose
	resultsChannel := make(chan error, 10)
	var wg sync.WaitGroup

	//
	// AAPL
	//
	if args.ProposeAAPL || args.ProposeAll {
		market := getMarket(markets, proposerVegawallet.PublicKey, MARKET_AAPL_MARKER)
		if market != nil {
			logger.Info("market already exist", zap.String("market", "AAPL"))
		} else {
			var settlementVegaAssetId string
			if args.FakeAsset {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "fUSDC")
			} else {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "tUSDC")
			}
			if err != nil {
				logger.Error("failed to get asset to propose AAPL", zap.Error(err))
			} else {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := proposeMarket(
						"AAPL", network.DataNodeClient, lastBlockData, proposerVegawallet, logger,
						proposals.NewAAPLMarketProposal(
							settlementVegaAssetId, 5, oraclePubKey,
							closingTime, enactmentTime,
							[]string{MARKET_AAPL_MARKER},
						),
					)
					if err != nil {
						resultsChannel <- err
					}
				}()
			}
		}
	}

	//
	// AAVEDAI
	//
	if args.ProposeAAVEDAI || args.ProposeAll {
		market := getMarket(markets, proposerVegawallet.PublicKey, MARKET_AAVEDAI_MARKER)
		if market != nil {
			logger.Info("market already exist", zap.String("market", "AAVEDAI"))
		} else {
			var settlementVegaAssetId string
			if args.FakeAsset {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "fDAI")
			} else {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "tDAI")
			}
			if err != nil {
				logger.Error("failed to get asset to propose AAVEDAI", zap.Error(err))
			} else {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := proposeMarket(
						"AAVEDAI", network.DataNodeClient, lastBlockData, proposerVegawallet, logger,
						proposals.NewAAVEDAIMarketProposal(
							settlementVegaAssetId, 5, oraclePubKey,
							closingTime, enactmentTime,
							[]string{MARKET_AAVEDAI_MARKER},
						),
					)
					if err != nil {
						resultsChannel <- err
					}
				}()
			}
		}
	}

	//
	// BTCUSD
	//
	if args.ProposeBTCUSD || args.ProposeAll {
		market := getMarket(markets, proposerVegawallet.PublicKey, MARKET_BTCUSD_MARKER)
		if market != nil {
			logger.Info("market already exist", zap.String("market", "BTCUSD"))
		} else {
			var settlementVegaAssetId string
			if args.FakeAsset {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "fDAI")
			} else {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "tDAI")
			}
			if err != nil {
				logger.Error("failed to get asset to propose BTCUSD", zap.Error(err))
			} else {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := proposeMarket(
						"BTCUSD", network.DataNodeClient, lastBlockData, proposerVegawallet, logger,
						proposals.NewBTCUSDMarketProposal(
							settlementVegaAssetId, 5, oraclePubKey,
							closingTime, enactmentTime,
							[]string{MARKET_BTCUSD_MARKER},
						),
					)
					if err != nil {
						resultsChannel <- err
					}
				}()
			}
		}
	}

	//
	// ETHBTC
	//
	if args.ProposeETHBTC || args.ProposeAll {
		market := getMarket(markets, proposerVegawallet.PublicKey, MARKET_ETHBTC_MARKER)
		if market != nil {
			logger.Info("market already exist", zap.String("market", "ETHBTC"))
		} else {
			var settlementVegaAssetId string
			if args.FakeAsset {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "fBTC")
			} else {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "tBTC")
			}
			if err != nil {
				logger.Error("failed to get asset to propose ETHBTC", zap.Error(err))
			} else {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := proposeMarket(
						"ETHBTC", network.DataNodeClient, lastBlockData, proposerVegawallet, logger,
						proposals.NewETHBTCMarketProposal(
							settlementVegaAssetId, 5, oraclePubKey,
							closingTime, enactmentTime,
							[]string{MARKET_ETHBTC_MARKER},
						),
					)
					if err != nil {
						resultsChannel <- err
					}
				}()
			}
		}
	}

	//
	// TSLA
	//
	if args.ProposeTSLA || args.ProposeAll {
		market := getMarket(markets, proposerVegawallet.PublicKey, MARKET_TSLA_MARKER)
		if market != nil {
			logger.Info("market already exist", zap.String("market", "TSLA"))
		} else {
			var settlementVegaAssetId string
			if args.FakeAsset {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "fEURO")
			} else {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "tEURO")
			}
			if err != nil {
				logger.Error("failed to get asset to propose TSLA", zap.Error(err))
			} else {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := proposeMarket(
						"TSLA", network.DataNodeClient, lastBlockData, proposerVegawallet, logger,
						proposals.NewTSLAMarketProposal(
							settlementVegaAssetId, 5, oraclePubKey,
							closingTime, enactmentTime,
							[]string{MARKET_TSLA_MARKER},
						),
					)
					if err != nil {
						resultsChannel <- err
					}
				}()
			}
		}
	}

	//
	// UNIDAI
	//
	if args.ProposeUNIDAI || args.ProposeAll {
		market := getMarket(markets, proposerVegawallet.PublicKey, MARKET_UNIDAI_MARKER)
		if market != nil {
			logger.Info("market already exist", zap.String("market", "UNIDAI"))
		} else {
			var settlementVegaAssetId string
			if args.FakeAsset {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "fEURO")
			} else {
				settlementVegaAssetId, err = getAssetIdForSymbol(assets, "tEURO")
			}
			if err != nil {
				logger.Error("failed to get asset to propose UNIDAI", zap.Error(err))
			} else {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := proposeMarket(
						"UNIDAI", network.DataNodeClient, lastBlockData, proposerVegawallet, logger,
						proposals.NewUNIDAIMarketProposal(
							settlementVegaAssetId, 5, oraclePubKey,
							closingTime, enactmentTime,
							[]string{MARKET_UNIDAI_MARKER},
						),
					)
					if err != nil {
						resultsChannel <- err
					}
				}()
			}
		}
	}

	wg.Wait()
	close(resultsChannel)

	for err := range resultsChannel {
		if err != nil {
			return fmt.Errorf("at least one proposal failed")
		}
	}

	return nil
}

func getAssetIdForSymbol(assets map[string]*vega.AssetDetails, symbol string) (string, error) {
	for id, asset := range assets {
		if asset.Symbol == symbol {
			return id, nil
		}
	}
	return "", fmt.Errorf("failed to find asset for symbol %s", symbol)
}

func getMarket(markets []*vega.Market, proposerPubKey string, metadataTag string) *vega.Market {
	for _, market := range markets {
		if slices.Contains(
			market.TradableInstrument.Instrument.GetFuture().OracleSpecForTradingTermination.PubKeys,
			proposerPubKey,
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

func proposeMarket(
	name string,
	dataNodeClient *vegaapi.DataNode,
	lastBlockData *vegaapipb.LastBlockHeightResponse,
	proposerVegawallet *wallet.VegaWallet,
	logger *zap.Logger,
	proposal *commandspb.ProposalSubmission,
) error {
	// Prepare vegawallet Transaction Request
	walletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_ProposalSubmission{
			ProposalSubmission: proposal,
		},
	}
	// Sign + Proof of Work vegawallet Transaction request
	signedTx, err := proposerVegawallet.SignTxWithPoW(&walletTxReq, lastBlockData)
	if err != nil {
		logger.Error("Failed to sign a trasnaction", zap.String("market", name),
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
	logger.Info("Submit transaction", zap.String("market", name),
		zap.String("proposer", proposerVegawallet.PublicKey))
	submitResponse, err := dataNodeClient.SubmitTransaction(submitReq)
	if err != nil {
		logger.Error("Failed to submit a trasnaction", zap.String("market", name),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", submitReq), zap.Error(err))
		return err
	}
	if !submitResponse.Success {
		logger.Error("Transaction submission response is not successful", zap.String("market", name),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", submitReq), zap.Any("response", submitResponse))
		return err
	}
	logger.Info("Successful Submision of Market Proposal", zap.String("market", name),
		zap.String("proposer", proposerVegawallet.PublicKey), zap.String("txHash", submitResponse.TxHash))
	return nil
}
