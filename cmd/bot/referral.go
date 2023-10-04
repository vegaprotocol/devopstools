package bot

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"

	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/bots"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/generate"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
)

type ReferralArgs struct {
	*BotArgs
}

var referralArgs ReferralArgs

// referralCmd represents the referral command
var referralCmd = &cobra.Command{
	Use:   "referral",
	Short: "manage bots in referral program",
	Long:  `manage bots in referral program`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunReferral(referralArgs); err != nil {
			referralArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	referralArgs.BotArgs = &botArgs

	BotCmd.AddCommand(referralCmd)
}

func RunReferral(args ReferralArgs) error {
	logger := args.Logger
	start := time.Now()
	logger.Info("Start referral", zap.Time("start", start))
	logger.Info("Connecting to nework")
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()
	logger.Info("Connected to network", zap.String("network", args.VegaNetworkName), zap.Duration("since start", time.Since(start)))

	logger.Info("Getting bots")
	botsAPIToken := args.BotsAPIToken
	if len(botsAPIToken) == 0 {
		botsAPIToken = network.BotsApiToken
	}

	traders, err := bots.GetResearchBots(args.VegaNetworkName, botsAPIToken)
	if err != nil {
		return err
	}
	logger.Info("Got bots", zap.Int("count", len(traders)), zap.Duration("since start", time.Since(start)))

	teamLeadersWallets := []*wallet.VegaWallet{}
	regularTradersWallets := []*wallet.VegaWallet{}

	for _, trader := range traders {
		wallet, err := trader.GetWallet()
		if err != nil {
			return fmt.Errorf("failed to get wallet for %s trader, %w", trader.PubKey, err)
		}
		if trader.WalletData.Index == 3 {
			teamLeadersWallets = append(teamLeadersWallets, wallet)
		} else {
			regularTradersWallets = append(regularTradersWallets, wallet)
		}
	}
	logger.Info("Divided bots into team leaders and regular wallets", zap.Int("leaders count", len(teamLeadersWallets)),
		zap.Int("regular wallet count", len(regularTradersWallets)), zap.Duration("since start", time.Since(start)))

	logger.Info("Create Referral Sets (teams)")
	if err := createReferralSetsForWallets(teamLeadersWallets, network, args.Logger); err != nil {
		return err
	}
	logger.Info("Created Referral Sets (teams)", zap.Duration("since start", time.Since(start)))

	logger.Info("Join Referral Sets (teams)")
	_ = regularTradersWallets
	logger.Info("Joined Referral Sets (teams)", zap.Duration("since start", time.Since(start)))

	// for _, trader := range traders {
	// 	_, err := trader.GetWallet()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if trader.WalletData.Index == 3 {
	// 		fmt.Printf("-----> %s\n", trader.PubKey)
	// 	} else {
	// 		fmt.Printf(" - %s (%d)\n", trader.PubKey, trader.WalletData.Index)
	// 	}
	// }
	// fmt.Printf("got wallets %s\n", time.Since(start))

	return nil
}

func createReferralSetsForWallets(
	referralTeamLeadVegawallet []*wallet.VegaWallet,
	network *veganetwork.VegaNetwork,
	logger *zap.Logger,
) error {
	var (
		minStake       *big.Int
		dataNodeClient = network.DataNodeClient
	)
	// get Referrals Tiers
	program, err := dataNodeClient.GetCurrentReferralProgram()
	if err != nil {
		return fmt.Errorf("failed to create referral sets, failed to get referral program, %w", err)
	}
	tiersMinStakeAmounts := make([]*big.Int, len(program.StakingTiers))
	for i, stakingTier := range program.StakingTiers {
		if stakeAmount, ok := new(big.Int).SetString(stakingTier.MinimumStakedTokens, 0); ok {
			tiersMinStakeAmounts[i] = ethutils.VegaTokenFromFullTokens(new(big.Float).SetInt(stakeAmount))
		} else {
			return fmt.Errorf("failed to convert %s to big.Int", stakingTier.MinimumStakedTokens)
		}
		if minStake == nil || minStake.Cmp(tiersMinStakeAmounts[i]) > 0 {
			minStake = tiersMinStakeAmounts[i]
		}
	}
	logger.Debug("Got Referral Program Staking Tiers", zap.Int("tiers count", len(tiersMinStakeAmounts)),
		zap.String("minStake", minStake.String()), zap.Any("tiers", tiersMinStakeAmounts))

	stakeByPubKey := map[string]*big.Int{}

	for _, wallet := range referralTeamLeadVegawallet {
		currentStake, err := dataNodeClient.GetPartyTotalStake(wallet.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to create referral sets, failed to get stake for %s, %w", wallet.PublicKey, err)
		}
		if currentStake.Cmp(minStake) < 0 {
			// rand stake
			rndIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(tiersMinStakeAmounts))))
			if err != nil {
				return err
			}
			expectedStake := tiersMinStakeAmounts[rndIdx.Int64()]
			stakeByPubKey[wallet.PublicKey] = expectedStake
			logger.Debug("Need to top up", zap.String("wallet", wallet.PublicKey),
				zap.String("current stake", currentStake.String()), zap.String("min stake", minStake.String()),
				zap.String("expected stake", expectedStake.String()))
		} else {
			stakeByPubKey[wallet.PublicKey] = currentStake
			logger.Debug("No need to top up", zap.String("wallet", wallet.PublicKey),
				zap.String("current stake", currentStake.String()), zap.String("min stake", minStake.String()))
		}
	}

	if err := Stake(stakeByPubKey, network, logger); err != nil {
		return err
	}

	return nil
}

func createReferralSet(
	creatorVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	errorMsg := fmt.Errorf("failed to create referral set by %s", creatorVegawallet.PrivateKey)
	teamName, err := generate.GenerateName()
	if err != nil {
		return fmt.Errorf("%s, %w", errorMsg, err)
	}
	teamName = fmt.Sprintf("Bots Team: %s", teamName)
	teamURL, err := generate.GenerateRandomWikiURL()
	if err != nil {
		return fmt.Errorf("%s, %w", errorMsg, err)
	}
	teamAvatar, err := generate.GenerateAvatarURL()
	if err != nil {
		return fmt.Errorf("%s, %w", errorMsg, err)
	}

	walletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: creatorVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_CreateReferralSet{
			CreateReferralSet: &commandspb.CreateReferralSet{
				IsTeam: true,
				Team: &commandspb.CreateReferralSet_Team{
					Name:      teamName,
					TeamUrl:   &teamURL,
					AvatarUrl: &teamAvatar,
				},
			},
		},
	}
	if err := governance.SubmitTx("vote on proposal", dataNodeClient, creatorVegawallet, logger, &walletTxReq); err != nil {
		return fmt.Errorf("%s, %w", errorMsg, err)
	}
	return nil
}

func Stake(
	expectedStakeForParty map[string]*big.Int,
	network *veganetwork.VegaNetwork,
	logger *zap.Logger,
) error {
	if len(expectedStakeForParty) == 0 {
		logger.Info("No part needs staking")
		return nil
	}
	//
	// Get missing stake
	//
	dataNodeClient := network.DataNodeClient
	totalMissingStake := big.NewInt(0)
	missingStakeByPubKey := map[string]*big.Int{}

	for pubKey, expectedStake := range expectedStakeForParty {
		currentStake, err := dataNodeClient.GetPartyTotalStake(pubKey)
		if err != nil {
			return fmt.Errorf("failed to create referral sets, failed to get stake for %s, %w", pubKey, err)
		}
		if currentStake.Cmp(expectedStake) < 0 {
			missingStake := new(big.Int).Sub(expectedStake, currentStake)
			missingStakeByPubKey[pubKey] = missingStake
			totalMissingStake = totalMissingStake.Add(totalMissingStake, missingStake)

			logger.Debug("Need to top up", zap.String("wallet", pubKey),
				zap.String("current stake", currentStake.String()),
				zap.String("expected stake", expectedStake.String()))
		} else {
			logger.Debug("No need to top up", zap.String("wallet", pubKey),
				zap.String("current stake", currentStake.String()),
				zap.String("expected stake", expectedStake.String()))
		}
	}
	if len(missingStakeByPubKey) == 0 {
		logger.Info("No part needs staking")
		return nil
	}
	//
	// Prepare source Ethereum wallet
	//
	vegaToken := network.SmartContracts.VegaToken
	wallet := network.NetworkMainWallet
	stakingBridge := network.SmartContracts.StakingBridge

	// TODO - check if there any pending transactions

	// check BALANCE
	currentBalance, err := vegaToken.BalanceOf(&bind.CallOpts{}, wallet.Address)
	if err != nil {
		return fmt.Errorf("failed to get balance, %w", err)
	}
	if currentBalance.Cmp(totalMissingStake) < 0 {
		return fmt.Errorf("wallet %s doesn't have enought tokens; current: %s, expected: %s",
			wallet.Address.Hex(), currentBalance.String(), totalMissingStake.String())
	}
	// check ALLOWANCE
	currentAllowance, err := vegaToken.Allowance(&bind.CallOpts{}, wallet.Address, stakingBridge.Address)
	if err != nil {
		return fmt.Errorf("failed to get allowance for staking bridge %s for party %s, %w",
			stakingBridge.Address.Hex(), wallet.Address, err)
	}
	if currentAllowance.Cmp(totalMissingStake) < 0 {
		// increase ALLOWANCE
		increaseAllowanceAmount := ethutils.VegaTokenFromFullTokens(big.NewFloat(1000000))
		if increaseAllowanceAmount.Cmp(totalMissingStake) < 0 {
			return fmt.Errorf("Want to stake more than 1kk tokens total - it is too much")
		}
		logger.Info("increasing allowance", zap.String("amount", increaseAllowanceAmount.String()),
			zap.String("ethWallet", wallet.Address.Hex()),
			zap.String("allowanceBefore", currentAllowance.String()),
			zap.String("tokenAddress", vegaToken.Address.Hex()))
		opts := wallet.GetTransactOpts()
		allowanceTx, err := vegaToken.IncreaseAllowance(opts, stakingBridge.Address, increaseAllowanceAmount)
		if err != nil {
			return fmt.Errorf("failed to increase allowance: %w", err)
		}
		// WAIT
		if err = ethutils.WaitForTransaction(network.EthClient, allowanceTx, time.Minute); err != nil {
			logger.Error("failed to increase allowance", zap.String("ethWallet", wallet.Address.Hex()),
				zap.String("tokenAddress", vegaToken.Address.Hex()), zap.Error(err))
			return fmt.Errorf("transaction failed to increase allowance: %w", err)
		}
		logger.Info("successfully increased allowance", zap.String("ethWallet", wallet.Address.Hex()),
			zap.String("tokenAddress", vegaToken.Address.Hex()))

	} else {
		logger.Info("no need to increase allowance", zap.String("tokenAddress", vegaToken.Address.Hex()),
			zap.String("ethWallet", wallet.Address.Hex()),
			zap.String("currentAllowance", currentAllowance.String()),
			zap.String("requiredAllowance", totalMissingStake.String()))
	}

	//
	// Stake
	//
	var (
		stakeTxs     = make(map[string]*ethTypes.Transaction, len(missingStakeByPubKey))
		failedCount  = 0
		successCount = 0
	)
	logger.Info("Sending stake transactions", zap.Int("count", len(missingStakeByPubKey)))
	for pubKey, missingStake := range missingStakeByPubKey {
		opts := wallet.GetTransactOpts()
		tx, err := stakingBridge.Stake(opts, missingStake, pubKey)
		if err != nil {
			failedCount += 1
			logger.Error("Failed to stake", zap.String("pub key", pubKey),
				zap.String("staking bridge", stakingBridge.Address.Hex()),
				zap.String("nonce", opts.Nonce.String()),
				zap.String("wallet", wallet.Address.Hex()), zap.Error(err))
		} else {
			stakeTxs[pubKey] = tx
			logger.Debug("Sent stake transaction", zap.String("pub key", pubKey),
				zap.String("staking bridge", stakingBridge.Address.Hex()),
				zap.String("nonce", opts.Nonce.String()),
				zap.String("wallet", wallet.Address.Hex()), zap.Any("tx", tx))
		}
	}
	logger.Info("Waiting for stake transactions", zap.Int("count", len(missingStakeByPubKey)),
		zap.Int("failed count", failedCount))
	//
	// wait for transactions to mint
	//
	for pubKey, tx := range stakeTxs {
		if tx == nil {
			continue
		}
		logger.Debug("waiting", zap.Any("tx", tx))
		if err = ethutils.WaitForTransaction(network.EthClient, tx, time.Minute); err != nil {
			failedCount += 1
			logger.Error("Stake transaction failed", zap.String("pub key", pubKey),
				zap.Uint64("nonce", tx.Nonce()), zap.Any("tx", tx), zap.Error(err))
		} else {
			successCount += 1
			logger.Debug("Stake transaction successful", zap.String("pub key", pubKey),
				zap.Uint64("nonce", tx.Nonce()), zap.String("tx", tx.Hash().Hex()))
		}
	}
	logger.Info("Finished staking", zap.Int("success count", successCount), zap.Int("failed count", failedCount),
		zap.Int("not needed count", len(expectedStakeForParty)-successCount-failedCount))
	if failedCount > 0 {
		return fmt.Errorf("stake for %d parties failed", failedCount)
	}
	return nil
}
