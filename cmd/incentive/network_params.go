package incentive

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"go.uber.org/zap"
)

type NetworkParamsArgs struct {
	*IncentiveArgs
	UpdateParams bool
}

var networkParamsArgs NetworkParamsArgs

// networkParamsCmd represents the networkParams command
var networkParamsCmd = &cobra.Command{
	Use:   "network-params",
	Short: "get or update (propose & vote) network params",
	Long:  `get or update (propose & vote) network params`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunNetworkParams(networkParamsArgs); err != nil {
			networkParamsArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	networkParamsArgs.IncentiveArgs = &incentiveArgs

	IncentiveCmd.AddCommand(networkParamsCmd)
	networkParamsCmd.PersistentFlags().BoolVar(&networkParamsArgs.UpdateParams, "update", false, "Update Network Parameter values with propose & vote")
}

type expectedNetworkParameter struct {
	Name          string
	ExpectedValue string
}

func expectedNetworkParams(env string) []expectedNetworkParameter {
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
		{Name: "rewards.activityStreak.benefitTiers", ExpectedValue: "{\"tiers\": [{\"minimum_activity_streak\": 1, \"reward_multiplier\": \"1.05\", \"vesting_multiplier\": \"1.05\"}, {\"minimum_activity_streak\": 24, \"reward_multiplier\": \"1.10\", \"vesting_multiplier\": \"1.10\"}, {\"minimum_activity_streak\": 48, \"reward_multiplier\": \"1.10\", \"vesting_multiplier\": \"1.15\"}, {\"minimum_activity_streak\": 62, \"reward_multiplier\": \"1.20\", \"vesting_multiplier\": \"1.20\"}]}"},
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
		{Name: "rewards.vesting.benefitTiers", ExpectedValue: `{"tiers": [{"minimum_quantum_balance": "1000", "reward_multiplier": "1.05"}]}`},
		{Name: "governance.proposal.transfer.minClose", ExpectedValue: "1m"},
		{Name: "governance.proposal.transfer.minEnact", ExpectedValue: "1m"},
	}

	if env == types.NetworkFairground {
		result = append(result, expectedNetworkParameter{
			Name:          "validators.epoch.length",
			ExpectedValue: "30m",
		})
	}

	return result
}

func RunNetworkParams(args NetworkParamsArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	toUpdate, err := checkNetworkParams(network)
	if err != nil {
		return err
	}

	if args.UpdateParams {
		updateCount, err := governance.ProposeAndVoteOnNetworkParameters(
			toUpdate, network.VegaTokenWhale, network.NetworkParams, network.DataNodeClient, args.Logger,
		)
		if err != nil {
			return err
		}
		if updateCount > 0 {
			if err := network.RefreshNetworkParams(); err != nil {
				return err
			}
		}
		_, err = checkNetworkParams(network)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkNetworkParams(network *veganetwork.VegaNetwork) (map[string]string, error) {
	yellowText := "\033[1;33m%s\033[0m"
	greenText := "\033[1;32m%s\033[0m"
	redText := "\033[1;31m%s\033[0m"
	toUpdate := map[string]string{}
	expectedParameters := expectedNetworkParams(network.Network)

	for _, param := range expectedParameters {
		fmt.Printf(" - %s = ", param.Name)
		if value, ok := network.NetworkParams.Params[param.Name]; ok {
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

	return toUpdate, nil
}
