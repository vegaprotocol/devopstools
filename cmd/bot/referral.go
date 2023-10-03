package bot

import (
	"fmt"
	"math/big"
	"os"
	"time"

	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/bots"
	"github.com/vegaprotocol/devopstools/generate"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/vegaapi"
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
	start := time.Now()
	fmt.Printf("start referral %s\n", start)
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()
	fmt.Printf("Connected to network %s\n", time.Since(start))

	botsAPIToken := args.BotsAPIToken
	if len(botsAPIToken) == 0 {
		botsAPIToken = network.BotsApiToken
	}

	traders, err := bots.GetResearchBots(args.VegaNetworkName, botsAPIToken)
	if err != nil {
		return err
	}
	fmt.Printf("got bots %s\n", time.Since(start))

	teamLeadersWallets := []*wallet.VegaWallet{}

	for _, trader := range traders {
		if trader.WalletData.Index == 3 {
			wallet, err := trader.GetWallet()
			if err != nil {
				return fmt.Errorf("failed to get wallet for %s trader, %w", trader.PubKey, err)
			}
			teamLeadersWallets = append(teamLeadersWallets, wallet)
		}
	}

	if err := createReferralSetsForWallets(teamLeadersWallets, network.DataNodeClient, args.Logger); err != nil {
		return err
	}

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
	creatorVegawallet []*wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	program, err := dataNodeClient.GetCurrentReferralProgram()
	if err != nil {
		return fmt.Errorf("failed to create referral sets, failed to get referral program, %w", err)
	}
	minStake, ok := new(big.Int).SetString(program.StakingTiers[0].MinimumStakedTokens, 0)
	if !ok {
		return fmt.Errorf("failed to convert %s to big.Int", program.StakingTiers[0].MinimumStakedTokens)
	}
	_ = minStake
	for _, wallet := range creatorVegawallet {
		stake, err := dataNodeClient.GetPartyTotalStake(wallet.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to create referral sets, failed to get stake for %s, %w", wallet.PublicKey, err)
		}
		fmt.Printf(" - %s: %d\n", wallet.PublicKey, stake)

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
