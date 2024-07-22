package bots

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"code.vegaprotocol.io/vega/core/netparams"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	v1 "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	eventspb "code.vegaprotocol.io/vega/protos/vega/events/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"
	"code.vegaprotocol.io/vega/wallet/wallet"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type AmmArgs struct {
	*Args

	timeout            time.Duration
	marketId           string
	poolSize           int
	vegaRecoveryPhrase string
}

const (
	ammSubmissionAmount     = 5000
	ammWhaleTopUpMultiplier = 5
	ammAgentTopUpMultiplier = 3
)

var ammArgs AmmArgs

var ammCmd = &cobra.Command{
	Use:   "amm",
	Short: "Setup the amm pool on given markets",
	Long:  "Setup the amm pool on given markets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setupAmm(ammArgs); err != nil {
			ammArgs.Logger.Error("Could not top up bots", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	ammArgs.Args = &args

	ammCmd.PersistentFlags().DurationVarP(&ammArgs.timeout, "timeout", "t", 15*time.Minute, "Timeout for the amm command")
	ammCmd.PersistentFlags().StringVarP(&ammArgs.marketId, "market-id", "m", "", "The market id")
	ammCmd.PersistentFlags().IntVarP(&ammArgs.poolSize, "size", "s", 20, "Timeout for the top up command")
	ammCmd.PersistentFlags().StringVarP(&ammArgs.vegaRecoveryPhrase, "vega-recovery-phrase", "p", "", "The vega recovery phrase for the amm pool submitter wallet")

	Cmd.AddCommand(ammCmd)
}

func setupAmm(args AmmArgs) error {
	ctx, cancelCommand := context.WithTimeout(context.Background(), topUpArgs.timeout)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(ctx, args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Info("Network file loaded", zap.String("name", cfg.Name.String()))

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Info("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Debug("Looking for healthy gRPC endpoints...")
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

	logger.Debug("Retrieving network parameters...")
	networkParams, err := datanodeClient.ListNetworkParameters(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}
	logger.Info("Network parameters retrieved")

	chainClients, err := ethereum.NewChainClients(ctx, cfg, networkParams, logger)
	if err != nil {
		return err
	}

	assets, err := datanodeClient.ListAssets(ctx)
	if err != nil {
		return fmt.Errorf("could not list assets from datanode: %w", err)
	}

	if len(assets) == 0 {
		return errors.New("no asset found on datanode")
	}

	activeMarkets, err := datanodeClient.GetAllMarketsWithState(ctx, datanode.ActiveMarkets)
	if err != nil {
		return fmt.Errorf("failed to get active markets")
	}

	var marketDetails *vegapb.Market
	for idx, market := range activeMarkets {
		if market.Id == args.marketId {
			marketDetails = activeMarkets[idx]
			break
		}
	}

	if marketDetails == nil {
		return fmt.Errorf("failed to find the %s market", args.marketId)
	}

	marketAssetsIds := datanode.AssetsIdsByMarket(marketDetails)
	if len(marketAssetsIds) == 0 {
		return fmt.Errorf("cannot obtain asset ids for given market")
	}

	submitterWallet, err := vega.LoadWallet("submitter", args.vegaRecoveryPhrase)
	if err != nil {
		return fmt.Errorf("failed to load submitter wallet: %w", err)
	}

	if err := vega.GenerateKeysUpToIndex(submitterWallet, uint32(args.poolSize+1)); err != nil {
		return fmt.Errorf("failed to generate %d keys from the submitter wallet: %w", args.poolSize, err)
	}

	topUpRegistry, err := determineAmmAgentsTopUpAmount(ctx, datanodeClient, assets, submitterWallet, args.poolSize, marketAssetsIds)
	if err != nil {
		return fmt.Errorf("failed to determine amm agents top up amount: %w", err)
	}

	whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	if err := vega.GenerateKeysUpToKey(whaleWallet, cfg.Network.Wallets.VegaTokenWhale.PublicKey); err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	whalePublicKey := cfg.Network.Wallets.VegaTokenWhale.PublicKey

	logger.Debug("Determining amounts to top up for whale...")
	whaleTopUpsByAsset, err := determineAmountsToTopUpForWhale(ctx, datanodeClient, whalePublicKey, assets, topUpRegistry, logger)
	if err != nil {
		return fmt.Errorf("failed to determine top up amounts for the whale: %w", err)
	}
	logger.Info("Amounts to top up for whale computed")

	if len(whaleTopUpsByAsset) == 0 {
		logger.Info("No top-up required for the whale")
	} else {
		logger.Debug("Depositing assets to whale...")
		if err := depositAssetsToWhale(ctx, whaleTopUpsByAsset, assets, datanodeClient, whalePublicKey, chainClients, logger); err != nil {
			return err
		}
		logger.Info("Whale has now enough funds to transfer to trading bots")
	}

	transfersToDo := countTransfersToDo(topUpRegistry)
	if networkParams.GetMaxTransfersPerEpoch() < int64(transfersToDo) {
		logger.Debug("Preparing network for transfers from whale to trading bots...")
		updateParams := map[string]string{
			netparams.TransferMaxCommandsPerEpoch: fmt.Sprintf("%d", transfersToDo*30),
		}
		if _, err := vega.UpdateNetworkParameters(ctx, whaleWallet, whalePublicKey, datanodeClient, updateParams, logger); err != nil {
			return fmt.Errorf("failed to prepare network for transfers: %w", err)
		}
		logger.Info("Network ready for transfers from whale to trading bots")
	}

	logger.Info("Transferring assets from whale to trading bots...")
	if err := transferAssetsFromWhaleToBots(ctx, datanodeClient, whaleWallet, whalePublicKey, topUpRegistry, logger); err != nil {
		return fmt.Errorf("failed to transfer assets from whale to one or more bots: %w", err)
	}
	logger.Info("Transfers done")

	marketData, err := datanodeClient.GetLatestMarketData(ctx, args.marketId)
	if err != nil {
		return fmt.Errorf("failed to get market data for market %s: %w", args.marketId, err)
	}

	submitterPubKeys := submitterWallet.ListPublicKeys()

	markPrice, _ := big.NewFloat(0).SetString(marketData.MarkPrice)
	upperBound, _ := big.NewFloat(0).Mul(
		markPrice,
		big.NewFloat(1.01),
	).Int(nil)
	upperBoundStr := upperBound.String()
	lowerBound, _ := big.NewFloat(0).Mul(
		markPrice,
		big.NewFloat(0.99),
	).Int(nil)
	lowerBoundStr := lowerBound.String()
	leverage := "10"

	// TODO add more assets when needed (for spot?)
	settlementAssetDetails := assets[marketAssetsIds[0]]
	for i := 0; i < args.poolSize; i++ {
		pubKey := submitterPubKeys[i].Key()

		activePartySubmissions, err := datanodeClient.ListAMMs(ctx, &pubKey, &args.marketId, true)
		if err != nil {
			return fmt.Errorf("failed to check active AMMs for party %s: %w", pubKey, err)
		}

		walletTxReq := walletpb.SubmitTransactionRequest{
			Command: &walletpb.SubmitTransactionRequest_SubmitAmm{
				SubmitAmm: &v1.SubmitAMM{
					MarketId: args.marketId,
					CommitmentAmount: big.NewInt(0).Mul(
						big.NewInt(ammSubmissionAmount),
						big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(settlementAssetDetails.Decimals)), nil),
					).String(),
					SlippageTolerance: "0.5",
					ProposedFee:       "0.001",
					ConcentratedLiquidityParameters: &v1.SubmitAMM_ConcentratedLiquidityParameters{
						UpperBound:           &upperBoundStr,
						LowerBound:           &lowerBoundStr,
						Base:                 marketData.MarkPrice,
						LeverageAtUpperBound: &leverage,
						LeverageAtLowerBound: &leverage,
					},
				},
			},
		}

		if len(activePartySubmissions) > 0 {
			walletTxReq = walletpb.SubmitTransactionRequest{
				Command: &walletpb.SubmitTransactionRequest_AmendAmm{
					AmendAmm: &v1.AmendAMM{
						MarketId:          args.marketId,
						CommitmentAmount:  nil,
						SlippageTolerance: "0.5",
						ProposedFee:       nil,
						ConcentratedLiquidityParameters: &v1.AmendAMM_ConcentratedLiquidityParameters{
							UpperBound:           &upperBoundStr,
							LowerBound:           &lowerBoundStr,
							Base:                 marketData.MarkPrice,
							LeverageAtUpperBound: &leverage,
							LeverageAtLowerBound: &leverage,
						},
					},
				},
			}
		}

		logger.Sugar().Infof("Submitting AMM for party %s", pubKey)
		resp, err := walletpkg.SendTransaction(ctx, submitterWallet, pubKey, &walletTxReq, datanodeClient)
		if err != nil {
			return fmt.Errorf("failed to submit amm with signature: %w", err)
		}

		if !resp.Success {
			logger.Sugar().Warnf("failed to submit amm: %s", resp.Data)
		}

		logger.Sugar().Infof("AMM Submission TX hash: %s", resp.TxHash)
	}

	// Wait for submissions to show up on te API and get their status
	logger.Info("Waiting for submission report")
	time.Sleep(30 * time.Second)

	for i := 0; i < args.poolSize; i++ {
		pubKey := submitterPubKeys[i].Key()

		amms, err := tools.RetryReturn(6, 5*time.Second, func() ([]*eventspb.AMM, error) {
			funcCtx, fCtxCancel := context.WithTimeout(ctx, 4*time.Second)
			defer fCtxCancel()

			amms, err := datanodeClient.ListAMMs(funcCtx, &pubKey, &args.marketId, true)
			if err != nil {
				return nil, fmt.Errorf("failed to list amm on report for party %s: %w", pubKey, err)
			}

			if len(amms) < 1 {
				return nil, fmt.Errorf("still waiting for active for party %s", pubKey)
			}

			return amms, nil
		})
		if err != nil {
			logger.Sugar().Errorf("no active submission found in api for party %s: %s", pubKey, err.Error())
		}

		if len(amms) < 1 {
			logger.Sugar().Warnf("Party %s: Submission is not active", pubKey)
		} else {
			logger.Sugar().Warnf("Party %s: AMM is active", pubKey)
		}
	}

	return nil
}

func determineAmmAgentsTopUpAmount(
	ctx context.Context,
	datanodeClient vegaapi.DataNodeClient,
	assets map[string]*vegapb.AssetDetails,
	submitterWallet wallet.Wallet,
	poolSize int,
	assetsIds []string,
) (map[string]AssetToTopUp, error) {
	result := map[string]AssetToTopUp{}

	ammAgentsPubKeys := submitterWallet.ListPublicKeys()
	if len(ammAgentsPubKeys) < poolSize {
		return nil, fmt.Errorf("not enough wallets generated for amm pool of %d", poolSize)
	}

	for i := 0; i <= poolSize; i++ {
		agentPubKey := ammAgentsPubKeys[i].Key()

		for _, assetId := range assetsIds {
			assetDetails, assetFound := assets[assetId]
			if !assetFound {
				return nil, fmt.Errorf("cannot find asset(%s) required for amm agent: %s", assetId, agentPubKey)
			}

			if assetDetails.GetErc20() == nil {
				return nil, fmt.Errorf("asset %s is not erc20 asset", assetId)
			}

			if _, assetResultFound := result[assetId]; !assetResultFound {
				result[assetId] = AssetToTopUp{
					Symbol:          assetDetails.Symbol,
					ContractAddress: assetDetails.GetErc20().ContractAddress,
					VegaAssetId:     assetId,
					TotalAmount:     types.NewAmount(assetDetails.Decimals),
					AmountsByParty:  map[string]*types.Amount{},
				}
			}

			currentAmmAgentGeneralAccount, err := datanodeClient.ListAccounts(ctx, agentPubKey, vegapb.AccountType_ACCOUNT_TYPE_GENERAL, &assetId)
			if err != nil {
				return nil, fmt.Errorf("failed to get general account balance for %s part", agentPubKey)
			}

			currentAmountWithoutZeros := big.NewInt(0)
			if len(currentAmmAgentGeneralAccount) > 0 {
				currentAmountWithoutZeros = big.NewInt(0).Div(
					currentAmmAgentGeneralAccount[0].Balance,
					big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(assetDetails.Decimals)), big.NewInt(0)),
				)
			}

			if currentAmountWithoutZeros.Cmp(big.NewInt(ammSubmissionAmount)) < 0 {
				result[assetId].TotalAmount.Add(big.NewFloat(0).Mul(big.NewFloat(ammAgentTopUpMultiplier), big.NewFloat(float64(ammSubmissionAmount)))) // Maybe we should add only missing amount, but for now meh...

				result[assetId].AmountsByParty[agentPubKey] = types.NewAmountFromMainUnit(
					big.NewFloat(0).Mul(big.NewFloat(ammAgentTopUpMultiplier), big.NewFloat(float64(ammSubmissionAmount))),
					assetDetails.Decimals,
				)
			}
		}
	}

	// multiplier for whale top up
	for _, assetId := range assetsIds {
		result[assetId].TotalAmount.Mul(big.NewFloat(float64(ammWhaleTopUpMultiplier)))
	}

	return result, nil
}
