package incentive

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type NetworkParamsArgs struct {
	*Args
	UpdateParams bool
	Change       bool
}

var networkParamsArgs NetworkParamsArgs

// networkParamsCmd represents the networkParams command
var networkParamsCmd = &cobra.Command{
	Use:   "network-params",
	Short: "Get or update (propose & vote) network params",
	Long:  "Get or update (propose & vote) network params",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunNetworkParams(networkParamsArgs); err != nil {
			networkParamsArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	networkParamsArgs.Args = &args

	Cmd.AddCommand(networkParamsCmd)
	networkParamsCmd.PersistentFlags().BoolVar(&networkParamsArgs.UpdateParams, "update", false, "Update network parameter values with propose & vote")
}

type expectedNetworkParameter struct {
	Name          string
	ExpectedValue string
}

func expectedNetworkParams() []expectedNetworkParameter {
	result := []expectedNetworkParameter{
		{Name: "rewards.vesting.baseRate", ExpectedValue: "0.0055"},
		{Name: "rewards.vesting.minimumTransfer", ExpectedValue: "10"},
		{Name: "referralProgram.maxReferralTiers", ExpectedValue: "10"},
		{Name: "referralProgram.maxReferralRewardFactor", ExpectedValue: "0.2"},
		{Name: "referralProgram.maxReferralDiscountFactor", ExpectedValue: "0.1"},
		{Name: "referralProgram.maxPartyNotionalVolumeByQuantumPerEpoch", ExpectedValue: "250000"},
		{Name: "referralProgram.minStakedVegaTokens", ExpectedValue: "0"},
		{Name: "referralProgram.maxReferralRewardProportion", ExpectedValue: "0.5"},
		{Name: "volumeDiscountProgram.maxBenefitTiers", ExpectedValue: "10"},
		{Name: "volumeDiscountProgram.maxVolumeDiscountFactor", ExpectedValue: "0.4"},
		{Name: "rewards.activityStreak.benefitTiers", ExpectedValue: "{\"tiers\": [{\"minimum_activity_streak\": 1, \"reward_multiplier\": \"1.05\", \"vesting_multiplier\": \"1.05\"}, {\"minimum_activity_streak\": 6, \"reward_multiplier\": \"1.10\", \"vesting_multiplier\": \"1.10\"}, {\"minimum_activity_streak\": 24, \"reward_multiplier\": \"1.10\", \"vesting_multiplier\": \"1.15\"}, {\"minimum_activity_streak\": 72, \"reward_multiplier\": \"1.20\", \"vesting_multiplier\": \"1.20\"}]}"},
		{Name: "rewards.activityStreak.inactivityLimit", ExpectedValue: "24"},
		{Name: "rewards.activityStreak.minQuantumOpenVolume", ExpectedValue: "100"},
		{Name: "rewards.activityStreak.minQuantumTradeVolume", ExpectedValue: "100"},
		{Name: "governance.proposal.referralProgram.minClose", ExpectedValue: "48h0m0s"},
		{Name: "governance.proposal.referralProgram.maxClose", ExpectedValue: "8760h0m0s"},
		{Name: "governance.proposal.referralProgram.minEnact", ExpectedValue: "48h0m0s"},
		{Name: "governance.proposal.referralProgram.maxEnact", ExpectedValue: "8760h0m0s"},
		{Name: "governance.proposal.referralProgram.requiredParticipation", ExpectedValue: "0.00001"},
		{Name: "governance.proposal.referralProgram.requiredMajority", ExpectedValue: "0.66"},
		{Name: "governance.proposal.referralProgram.minProposerBalance", ExpectedValue: "1"},
		{Name: "governance.proposal.referralProgram.minVoterBalance", ExpectedValue: "1"},
		{Name: "governance.proposal.VolumeDiscountProgram.minClose", ExpectedValue: "48h0m0s"},
		{Name: "governance.proposal.VolumeDiscountProgram.maxClose", ExpectedValue: "8760h0m0s"},
		{Name: "spam.protection.max.createReferralSet", ExpectedValue: "3"},
		{Name: "spam.protection.max.updateReferralSet", ExpectedValue: "3"},
		{Name: "spam.protection.max.applyReferralCode", ExpectedValue: "5"},
		{Name: "governance.proposal.VolumeDiscountProgram.minEnact", ExpectedValue: "48h0m0s"},
		{Name: "governance.proposal.VolumeDiscountProgram.maxEnact", ExpectedValue: "8760h0m0s"},
		{Name: "governance.proposal.VolumeDiscountProgram.requiredParticipation", ExpectedValue: "0.00001"},
		{Name: "governance.proposal.VolumeDiscountProgram.requiredMajority", ExpectedValue: "0.66"},
		{Name: "governance.proposal.VolumeDiscountProgram.minProposerBalance", ExpectedValue: "1"},
		{Name: "governance.proposal.VolumeDiscountProgram.minVoterBalance", ExpectedValue: "1"},
		{Name: "rewards.vesting.benefitTiers", ExpectedValue: `{"tiers": [{"minimum_quantum_balance": "10", "reward_multiplier": "1.05"}, {"minimum_quantum_balance": "100", "reward_multiplier": "1.10"},{"minimum_quantum_balance": "1000", "reward_multiplier": "1.15"},{"minimum_quantum_balance": "10000", "reward_multiplier": "1.20"}]}`},
		{Name: "governance.proposal.transfer.minClose", ExpectedValue: "1m"},
		{Name: "governance.proposal.transfer.minEnact", ExpectedValue: "1m"},
		{Name: "governance.proposal.updateAsset.minClose", ExpectedValue: "5s"},
		{Name: "governance.proposal.updateAsset.minEnact", ExpectedValue: "5s"},
		{Name: "ethereum.oracles.enabled", ExpectedValue: "1"},
		{Name: "validators.epoch.length", ExpectedValue: "2h"},
		{Name: "auction.LongBlock", ExpectedValue: `{"threshold_and_duration": [{"threshold":"10s","duration":"1m"}]}`},
		{Name: "limits.markets.ammPoolEnabled", ExpectedValue: "1"},
		{Name: "market.amm.minCommitmentQuantum", ExpectedValue: "1"},
		{Name: "market.liquidity.maxAmmCalculationLevels", ExpectedValue: "100"},
	}

	return result
}

func RunNetworkParams(args NetworkParamsArgs) error {
	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(ctx, args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Debug("Network file loaded", zap.String("name", cfg.Name.String()))

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Debug("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Debug("Looking for healthy gRPC endpoints...")
	healthyEndpoints := tools.FilterHealthyGRPCEndpoints(endpoints)
	if len(healthyEndpoints) == 0 {
		return fmt.Errorf("no healthy gRPC endpoint found on configured datanodes")
	}
	logger.Debug("Healthy gRPC endpoints found", zap.Strings("endpoints", healthyEndpoints))

	datanodeClient := datanode.New(healthyEndpoints, 3*time.Second, args.Logger.Named("datanode"))

	logger.Debug("Connecting to a datanode's gRPC endpoint...")
	dialCtx, cancelDialing := context.WithTimeout(ctx, 2*time.Second)
	defer cancelDialing()
	datanodeClient.MustDialConnection(dialCtx) // blocking
	logger.Debug("Connected to a datanode's gRPC node", zap.String("node", datanodeClient.Target()))

	logger.Debug("Retrieving network parameters...")
	networkParams, err := datanodeClient.ListNetworkParameters(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}
	logger.Debug("Network parameters retrieved")

	whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}
	publicKey := cfg.Network.Wallets.VegaTokenWhale.PublicKey
	if err := vega.GenerateKeysUpToKey(whaleWallet, publicKey); err != nil {
		return fmt.Errorf("could not generate whale keys: %w", err)
	}

	toUpdate := checkNetworkParams(networkParams)

	if args.UpdateParams {
		updateCount, err := governance.ProposeAndVoteOnNetworkParameters(ctx, toUpdate, whaleWallet, publicKey, networkParams, datanodeClient, logger)
		if err != nil {
			return err
		}
		if updateCount > 0 {
			logger.Debug("Retrieving network parameters...")
			updatedNetworkParams, err := datanodeClient.ListNetworkParameters(ctx)
			if err != nil {
				return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
			}
			logger.Debug("Network parameters retrieved")
			networkParams = updatedNetworkParams
		}
		_ = checkNetworkParams(networkParams)
	}

	return nil
}

func checkNetworkParams(netParams *types.NetworkParams) map[string]string {
	yellowText := "\033[1;33m%s\033[0m"
	greenText := "\033[1;32m%s\033[0m"
	redText := "\033[1;31m%s\033[0m"
	toUpdate := map[string]string{}
	expectedParameters := expectedNetworkParams()

	for _, param := range expectedParameters {
		fmt.Printf(" - %s = ", param.Name)
		if value, ok := netParams.Params[param.Name]; ok {
			if value == param.ExpectedValue {
				fmt.Printf(greenText, fmt.Sprintf("%s (ok)\n", value))
			} else {
				fmt.Printf(redText, fmt.Sprintf("%s != expected %s\n", value, param.ExpectedValue))
				toUpdate[param.Name] = param.ExpectedValue
			}
		} else {
			fmt.Printf(yellowText, "does not exist\n")
		}
	}

	return toUpdate
}
