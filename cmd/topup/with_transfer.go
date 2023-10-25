package topup

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	vegacmd "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	v1 "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/bots"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
)

const (
	DefaultTopUpValue = 10000
	TraderTopUpFactor = 3.0
	WhaleTopUpFactor  = 10.0
)

type TopUpWithTransferArgs struct {
	*TopUpArgs
	VegaNetworkName string
	TradersURL      string
}

var topUpWithTransferArgs TopUpWithTransferArgs

// topUpWithTransferCmd represents the traderbot command
var topUpWithTransferCmd = &cobra.Command{
	Use:   "with-transfer",
	Short: "TopUp parties on network with vega transfer",
	Long:  `TopUp parties on network with vega transfer`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunTopUpWithTransfer(topUpWithTransferArgs); err != nil {
			traderbotArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	topUpWithTransferArgs.TopUpArgs = &topUpArgs

	TopUpCmd.AddCommand(topUpWithTransferCmd)
	topUpWithTransferCmd.PersistentFlags().StringVar(&topUpWithTransferArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := topUpWithTransferCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	topUpWithTransferCmd.PersistentFlags().StringVar(&topUpWithTransferArgs.TradersURL, "traders-url", "", "Traders URL")
}

type AssetTopUp struct {
	Symbol          string
	ContractAddress string
	VegaAssetId     string
	TotalAmount     *big.Float
	Parties         map[string]*big.Float
}

func RunTopUpWithTransfer(args TopUpWithTransferArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return fmt.Errorf("failed to create vega network object: %w", err)
	}
	defer network.Disconnect()
	networkAssets, err := network.DataNodeClient.GetAssets()
	if err != nil {
		return fmt.Errorf("failed to get assets from datanode: %w", err)
	}

	traders, err := tools.RetryReturn(10, 5*time.Second, func() (map[string]bots.BotTraders, error) {
		return bots.GetBotTradersWithURL(args.VegaNetworkName, args.TradersURL)
	})
	if err != nil {
		return fmt.Errorf("failed to read traders: %w", err)
	}

	tradersTopUpRegistry, err := determineTradersTopUpAmount(args.Logger, networkAssets, traders)
	if err != nil {
		return fmt.Errorf("failed to determine top up amounts for the assets: %w", err)
	}
	args.Logger.Info("")

	whaleTopUpRegistry, err := determineWhaleTopUpAmount(
		args.Logger,
		network.DataNodeClient,
		networkAssets,
		tradersTopUpRegistry,
		network.VegaTokenWhale.PublicKey,
	)
	if err != nil {
		return fmt.Errorf("failed to compute top up amount for whale: %w", err)
	}

	for assetId, amount := range whaleTopUpRegistry {
		fmt.Printf("whale top up %s: %s", assetId, tools.AsIntStringFromFloat(amount))

		if err := depositToWhale(
			args.Logger,
			network,
			network.VegaTokenWhale.PublicKey,
			networkAssets[assetId],
			amount,
		); err != nil {
			return fmt.Errorf("failed to deposit %s: %w", assetId, err)
		}
	}

	args.Logger.Info("Waiting for money to appear on the whale wallet")

	args.Logger.Info("Voting network parameters")
	numberOfTransfers := countTransferNumbers(tradersTopUpRegistry)
	if err := prepareNetworkForTransfer(args.Logger, network, numberOfTransfers); err != nil {
		return fmt.Errorf("failed to prepare network: %w", err)
	}

	// wait for money
	for assetId, amount := range whaleTopUpRegistry {
		assetDetails := networkAssets[assetId]
		if err := tools.RetryRun(15, 6*time.Second, func() error {
			return waitForMoney(
				network.DataNodeClient,
				network.VegaTokenWhale.PublicKey,
				assetId,
				assetDetails,
				amount,
			)
		}); err != nil {
			return fmt.Errorf(
				"failed to wait for money on whale wallet: waiting for money to appear on whale account for token %s timed out",
				assetDetails.Name,
			)
		}

		args.Logger.Info("Whale has enough tokens", zap.String("asset", assetDetails.Name))
	}

	stats, err := transferMoneyFromWhaleToBots(
		args.Logger,
		network,
		network.VegaTokenWhale,
		tradersTopUpRegistry,
	)

	args.Logger.Info(
		"Transfer of money from whale to bots finished",
		zap.Int("successful", stats.successful),
		zap.Int("failed", stats.failed),
		zap.Int("total", stats.failed+stats.successful),
		zap.String("failed-transactions", strings.Join(stats.failedTransactions, ", ")),
		zap.String("successful-transactions", strings.Join(stats.successfulTransactions, ", ")),
	)

	if err != nil {
		return fmt.Errorf("failed to transfer money from whale to bots for one or more parties: %w", err)
	}

	return nil
}

func countTransferNumbers(topUpRegistry map[string]AssetTopUp) int {
	transferNumbers := 0

	for _, entry := range topUpRegistry {
		transferNumbers = transferNumbers + len(entry.Parties)
	}

	return transferNumbers + (10 * transferNumbers / 100)
}

func prepareNetworkForTransfer(logger *zap.Logger, network *veganetwork.VegaNetwork, numberOfBots int) error {
	updateParams := map[string]string{
		"spam.protection.maxUserTransfersPerEpoch": fmt.Sprintf("%d", numberOfBots),
	}

	logger.Info("Refreshing network parameters")
	if err := network.RefreshNetworkParams(); err != nil {
		return fmt.Errorf("failed to refresh network parameters: %w", err)
	}

	logger.Info("Getting value for spam.protection.maxUserTransfersPerEpoch network parameter")
	maxUserTransfersPerEpoch, paramExists := network.NetworkParams.Params["spam.protection.maxUserTransfersPerEpoch"]
	if !paramExists {
		return fmt.Errorf("failed to get spam.protection.maxUserTransfersPerEpoch value from the network parameters: parameter does not exist")
	}

	maxUserTransfersPerEpochInt, err := strconv.Atoi(maxUserTransfersPerEpoch)
	if err != nil {
		return fmt.Errorf("failed to convert value of spam.protection.maxUserTransfersPerEpoch to int: %w", err)
	}

	if numberOfBots < maxUserTransfersPerEpochInt {
		logger.Sugar().Infof(
			"spam.protection.maxUserTransfersPerEpoch network parameter does not need to be modified. Current value: %d, expected value at least %d",
			maxUserTransfersPerEpochInt,
			numberOfBots,
		)
		return nil
	}

	updateCount, err := governance.ProposeAndVoteOnNetworkParameters(
		updateParams, network.VegaTokenWhale, network.NetworkParams, network.DataNodeClient, logger,
	)

	if err != nil {
		return fmt.Errorf("failed to propose and vote on network parameters: %w", err)
	}
	if updateCount > 0 {
		if err := network.RefreshNetworkParams(); err != nil {
			return fmt.Errorf("failed to refresh network parameters: %w", err)
		}
	}
	for name, expectedValue := range updateParams {
		if network.NetworkParams.Params[name] != expectedValue {
			return fmt.Errorf("failed to update Network Parameter '%s', current value: '%s', expected value: '%s'",
				name, network.NetworkParams.Params[name], expectedValue,
			)
		}
	}
	return nil
}

func waitForMoney(
	dataNodeClient vegaapi.DataNodeClient,
	partyId, vegaAssetId string,
	assetDetails *vega.AssetDetails,
	requiredMoney *big.Float,
) error {
	whaleFund, err := dataNodeClient.GetFunds(partyId, vega.AccountType_ACCOUNT_TYPE_GENERAL, &vegaAssetId)
	if err != nil {
		return fmt.Errorf("failed to get funds for whale(%s): %w", partyId, err)
	}

	requiredFunds, _ := requiredMoney.Int64()
	requiredFundsWithZeros := big.NewInt(0).
		Mul(
			big.NewInt(requiredFunds),
			big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(assetDetails.Decimals)), nil),
		)
	whaleFundsWithZeros := big.NewInt(0)
	if len(whaleFund) > 0 {
		whaleFundsWithZeros = whaleFund[0].Balance
	}

	if whaleFundsWithZeros.Cmp(requiredFundsWithZeros) < 0 {
		return fmt.Errorf("whale wallet does not have enough tokens")
	}

	return nil
}

type transferStats struct {
	successful             int
	successfulTransactions []string

	failed             int
	failedTransactions []string
}

func transferMoneyFromWhaleToBots(
	logger *zap.Logger,
	network *veganetwork.VegaNetwork,
	whaleWallet *wallet.VegaWallet,
	registry map[string]AssetTopUp,
) (*transferStats, error) {
	networkAssets, err := network.DataNodeClient.GetAssets()
	if err != nil {
		return nil, fmt.Errorf("failed to get assets from network: %w", err)
	}

	stats := &transferStats{}
	var result *multierror.Error

	for asset, entry := range registry {
		assetMultiplier := big.NewInt(0).
			Exp(
				big.NewInt(10),
				big.NewInt(0).SetUint64(networkAssets[asset].Decimals),
				nil,
			).
			Int64()

		for receiverPartyId, amount := range entry.Parties {
			lastBlockData, err := network.DataNodeClient.LastBlockData()
			if err != nil {
				logger.Error(
					"failed to get statistics before transfer transaction is signed",
					zap.String("token", entry.Symbol),
					zap.String("receiver", receiverPartyId),
					zap.String("amount", tools.AsIntStringFromFloat(amount)),
					zap.Error(err),
				)
				multierror.Append(result, fmt.Errorf(
					"failed to get statistics before transfer transaction is signed: %w",
					err,
				))

				stats.failed = stats.failed + 1
				continue
			}

			transferAmountWithZeros := big.NewFloat(0).
				Mul(
					amount,
					big.NewFloat(0).SetInt64(assetMultiplier),
				)
			signedTransaction, err := whaleWallet.SignTxWithPoW(&v1.SubmitTransactionRequest{
				PubKey: whaleWallet.PublicKey,
				Command: &v1.SubmitTransactionRequest_Transfer{
					Transfer: &vegacmd.Transfer{
						Reference:       fmt.Sprintf("Transfer from whale to %s", receiverPartyId),
						FromAccountType: vega.AccountType_ACCOUNT_TYPE_GENERAL,
						ToAccountType:   vega.AccountType_ACCOUNT_TYPE_GENERAL,
						To:              receiverPartyId,
						Asset:           asset,
						Amount:          tools.AsIntStringFromFloat(transferAmountWithZeros),
						Kind: &vegacmd.Transfer_OneOff{
							OneOff: &vegacmd.OneOffTransfer{
								DeliverOn: 0,
							},
						},
					},
				},
			}, lastBlockData)

			if err != nil {
				logger.Error(
					"failed to sign transaction with PoW",
					zap.String("token", entry.Symbol),
					zap.String("receiver", receiverPartyId),
					zap.String("amount", tools.AsIntStringFromFloat(amount)),
					zap.Error(err),
				)

				multierror.Append(result, fmt.Errorf(
					"failed to sign the %s transfer for %s bot transaction with vega wallet: %w",
					entry.Symbol,
					receiverPartyId,
					err,
				))

				stats.failed = stats.failed + 1
				continue
			}

			resp, err := network.DataNodeClient.SubmitTransaction(&vegaapipb.SubmitTransactionRequest{
				Tx:   signedTransaction,
				Type: vegaapipb.SubmitTransactionRequest_TYPE_SYNC,
			})

			if err != nil {
				logger.Error(
					"failed to send the signed transaction",
					zap.String("token", entry.Symbol),
					zap.String("receiver", receiverPartyId),
					zap.String("amount", tools.AsIntStringFromFloat(amount)),
					zap.Error(err),
				)

				multierror.Append(result, fmt.Errorf(
					"failed to send the %s transfer for %s bot: %w",
					entry.Symbol,
					receiverPartyId,
					err,
				))

				stats.failed = stats.failed + 1
				continue
			}

			if !resp.Success {
				logger.Error(
					"Sent transaction is not successful",
					zap.String("token", entry.Symbol),
					zap.String("receiver", receiverPartyId),
					zap.String("amount", tools.AsIntStringFromFloat(amount)),
					zap.String("amount-with-zeros", tools.AsIntStringFromFloat(transferAmountWithZeros)),
					zap.String("reason", resp.Data),
				)

				multierror.Append(result,
					fmt.Errorf("failed to successfully send transfer transaction to the network: %s", resp.Data),
				)

				stats.failed = stats.failed + 1
				stats.failedTransactions = append(stats.failedTransactions, resp.TxHash)
				continue
			}

			logger.Info(
				"Tokens has been sent from whale to bot wallet",
				zap.String("token", entry.Symbol),
				zap.String("amount", tools.AsIntStringFromFloat(amount)),
				zap.String("amount-with-zeros", tools.AsIntStringFromFloat(transferAmountWithZeros)),
				zap.String("receiver", receiverPartyId),
				zap.String("transaction", resp.TxHash),
			)
			stats.successful = stats.successful + 1
			stats.successfulTransactions = append(stats.successfulTransactions, resp.TxHash)

			time.Sleep(100 * time.Millisecond)
		}
	}

	return stats, result.ErrorOrNil()
}

func depositToWhale(
	logger *zap.Logger,
	network *veganetwork.VegaNetwork,
	partyId string,
	asset *vega.AssetDetails,
	amount *big.Float,
) error {
	if asset.GetErc20() == nil {
		return fmt.Errorf("token %s is not an erc20 token", asset.Symbol)
	}

	if err := depositERC20TokenToParties(
		network,
		asset.GetErc20().ContractAddress,
		[]string{partyId},
		amount,
		logger,
	); err != nil {
		return fmt.Errorf("failed to deposit erc20 token: %w", err)
	}

	return nil
}

func determineWhaleTopUpAmount(
	logger *zap.Logger,
	dataNodeClient vegaapi.DataNodeClient,
	assets map[string]*vega.AssetDetails,
	tradersRegistry map[string]AssetTopUp,
	whalePartyId string,
) (map[string]*big.Float, error) {
	result := map[string]*big.Float{}

	for _, traderRegistryEntry := range tradersRegistry {
		assetDetails, assetExists := assets[traderRegistryEntry.VegaAssetId]
		if !assetExists {
			return nil, fmt.Errorf(
				"failed to find asset on network: whale needs to topup the %s asset but it does not exist on the network",
				traderRegistryEntry.VegaAssetId,
			)
		}

		whaleFund, err := dataNodeClient.GetFunds(whalePartyId, vega.AccountType_ACCOUNT_TYPE_GENERAL, &traderRegistryEntry.VegaAssetId)
		if err != nil {
			return nil, fmt.Errorf("failed to get funds for whale(%s): %w", whalePartyId, err)
		}

		requiredFunds, _ := traderRegistryEntry.TotalAmount.Int64()
		requiredFundsWithZeros := big.NewInt(0).
			Mul(
				big.NewInt(requiredFunds),
				big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(assetDetails.Decimals)), nil),
			)
		whaleFundsWithZeros := big.NewInt(0)
		if len(whaleFund) > 0 {
			whaleFundsWithZeros = whaleFund[0].Balance
		}
		whaleFunds := big.NewInt(0).
			Div(
				whaleFundsWithZeros,
				big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(assetDetails.Decimals)), nil),
			)

		if whaleFundsWithZeros.Cmp(requiredFundsWithZeros) > -1 {
			logger.Info(
				fmt.Sprintf(
					"Whale does not need top up for the %s asset. It already has enough funds",
					traderRegistryEntry.Symbol,
				),
				zap.Int64("Required funds", requiredFunds),
				zap.String("Wallet funds", whaleFunds.String()),
				zap.String("Party", whalePartyId),
			)
			continue
		}

		topUpAmountNonZeros := big.NewFloat(0).
			Mul(
				traderRegistryEntry.TotalAmount,
				big.NewFloat(WhaleTopUpFactor),
			)

		logger.Info(
			fmt.Sprintf(
				"Whale need top up for the %s asset",
				traderRegistryEntry.Symbol,
			),
			zap.Int64("Required funds", requiredFunds),
			zap.String("Wallet funds", whaleFunds.String()),
			zap.String("Top up amount", tools.AsIntStringFromFloat(topUpAmountNonZeros)),
			zap.String("Party", whalePartyId),
		)

		result[traderRegistryEntry.VegaAssetId] = topUpAmountNonZeros
	}

	return result, nil
}

func determineTradersTopUpAmount(
	logger *zap.Logger,
	assets map[string]*vega.AssetDetails,
	traders map[string]bots.BotTraders,
) (map[string]AssetTopUp, error) {
	topUpRegistry := map[string]AssetTopUp{}

	for _, traderDetails := range traders {
		assetDetails, assetExists := assets[traderDetails.Parameters.SettlementVegaAssetID]
		if !assetExists {
			return nil, fmt.Errorf(
				"failed to find asset on network: bot is using the %s asset but it does not exist on the network",
				traderDetails.Parameters.SettlementVegaAssetID,
			)
		}

		if _, registryExists := topUpRegistry[traderDetails.Parameters.SettlementVegaAssetID]; !registryExists {
			topUpRegistry[traderDetails.Parameters.SettlementVegaAssetID] = AssetTopUp{
				Symbol:          assetDetails.Symbol,
				ContractAddress: traderDetails.Parameters.SettlementEthereumContractAddress,
				VegaAssetId:     traderDetails.Parameters.SettlementVegaAssetID,
				TotalAmount:     big.NewFloat(0),
				Parties:         map[string]*big.Float{},
			}
		}

		currentEntry := topUpRegistry[traderDetails.Parameters.SettlementVegaAssetID]

		requiredTopUp := necessaryTopUp(
			traderDetails.Parameters.CurrentBalance,
			traderDetails.Parameters.WantedTokens,
			TraderTopUpFactor,
		)

		if requiredTopUp == 0 {
			continue
		}

		currentEntry.Parties[traderDetails.PubKey] = big.NewFloat(requiredTopUp)
		currentEntry.TotalAmount = big.NewFloat(0.0).Add(currentEntry.TotalAmount, big.NewFloat(requiredTopUp))

		logger.Info(
			"Required top up for party",
			zap.String("party-id", traderDetails.PubKey),
			zap.Float64("amount", requiredTopUp),
			zap.String("asset", assetDetails.Name),
		)

		topUpRegistry[traderDetails.Parameters.SettlementVegaAssetID] = currentEntry
	}

	return topUpRegistry, nil
}

func necessaryTopUp(currentBalance, wantedBalance, factor float64) float64 {
	if wantedBalance < 0.01 || wantedBalance > currentBalance {
		return wantedBalance * factor
	}

	// top up not required
	return 0
}
