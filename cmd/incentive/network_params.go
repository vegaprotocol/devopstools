package incentive

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"go.uber.org/zap"
)

type NetworkParamsArgs struct {
	*IncentiveArgs
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
}

var expectedNetworkParams = []struct {
	Name          string
	ExpectedValue string
}{
	{Name: "rewards.vesting.baseRate", ExpectedValue: "0.0055"},
	{Name: "rewards.vesting.minimumTransfer", ExpectedValue: "10"},
	{Name: "referralProgram.maxReferralTiers", ExpectedValue: "100"},
	{Name: "referralProgram.maxStakingTiers", ExpectedValue: "0"},
	{Name: "referralProgram.maxReferralRewardFactor", ExpectedValue: "0.2"},
	{Name: "referralProgram.maxReferralDiscountFactor", ExpectedValue: "0.1"},
	{Name: "referralProgram.maxPartyNotionalVolumeByQuantumPerEpoch", ExpectedValue: "250000"},
	{Name: "referralProgram.minStakedVegaTokens", ExpectedValue: "0"},
	{Name: "referralProgram.maxReferralRewardProportion", ExpectedValue: "0.5"},
	{Name: "volumeDiscountProgram.maxBenefitTiers", ExpectedValue: "100"},
	{Name: "volumeDiscountProgram.maxVolumeDiscountFactor", ExpectedValue: "0.4"},
	{Name: "rewards.activityStreak.benefitTiers", ExpectedValue: "(See next page)"},
	{Name: "rewards.activityStreak.inactivityLimit", ExpectedValue: "24"},
	{Name: "rewards.activityStreak.minQuantumOpenVolume", ExpectedValue: "100"},
	{Name: "rewards.activityStreak.minQuantumnTradeVolume", ExpectedValue: "100"},
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
}

var expectedNetworkParamsOld = map[string]string{
	"rewards.vesting.baseRate":                                "0.0055",
	"rewards.vesting.minimumTransfer":                         "10",
	"referralProgram.maxReferralTiers":                        "100",
	"referralProgram.maxStakingTiers":                         "0",
	"referralProgram.maxReferralRewardFactor":                 "0.2",
	"referralProgram.maxReferralDiscountFactor":               "0.1",
	"referralProgram.maxPartyNotionalVolumeByQuantumPerEpoch": "250000",
	"referralProgram.minStakedVegaTokens":                     "0",
	"referralProgram.maxReferralRewardProportion":             "0.5",
	"volumeDiscountProgram.maxBenefitTiers":                   "100",
	"volumeDiscountProgram.maxVolumeDiscountFactor":           "0.4",
	//"rewards.activityStreak.benefitTiers":                             "(See next page)",
	"rewards.activityStreak.inactivityLimit":                          "24",
	"rewards.activityStreak.minQuantumOpenVolume":                     "100",
	"rewards.activityStreak.minQuantumnTradeVolume":                   "100",
	"governance.proposal.referralProgram.minClose":                    "48h0m0s",
	"governance.proposal.referralProgram.maxClose":                    "8760h0m0s",
	"governance.proposal.referralProgram.minEnact":                    "48h0m0s",
	"governance.proposal.referralProgram.maxEnact":                    "8760h0m0s",
	"governance.proposal.referralProgram.requiredParticipation":       "0.00001",
	"governance.proposal.referralProgram.requiredMajority":            "0.66",
	"governance.proposal.referralProgram.minProposerBalance":          "1",
	"governance.proposal.referralProgram.minVoterBalance":             "1",
	"governance.proposal.VolumeDiscountProgram.minClose":              "48h0m0s",
	"governance.proposal.VolumeDiscountProgram.maxClose":              "8760h0m0s",
	"spam.protection.max.createReferralSet":                           "3",
	"spam.protection.max.updateReferralSet":                           "3",
	"spam.protection.max.applyReferralCode":                           "5",
	"governance.proposal.VolumeDiscountProgram.minEnact":              "48h0m0s",
	"governance.proposal.VolumeDiscountProgram.maxEnact":              "8760h0m0s",
	"governance.proposal.VolumeDiscountProgram.requiredParticipation": "0.00001",
	"governance.proposal.VolumeDiscountProgram.requiredMajority":      "0.66",
	"governance.proposal.VolumeDiscountProgram.minProposerBalance":    "1",
	"governance.proposal.VolumeDiscountProgram.minVoterBalance":       "1",
}

func RunNetworkParams(args NetworkParamsArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	if err := checkNetworkParams(network); err != nil {
		return err
	}

	return nil
}

func checkNetworkParams(network *veganetwork.VegaNetwork) error {
	yellowText := "\033[1;33m%s\033[0m"
	greenText := "\033[1;32m%s\033[0m"
	redText := "\033[1;31m%s\033[0m"
	for _, param := range expectedNetworkParams {
		fmt.Printf(" - %s = ", param.Name)
		if value, ok := network.NetworkParams.Params[param.Name]; ok {
			if value == param.ExpectedValue {
				fmt.Printf(greenText, fmt.Sprintf("%s (ok)\n", value))
			} else {
				fmt.Printf(redText, fmt.Sprintf("%s != expected %s\n", value, param.ExpectedValue))
			}
		} else {
			fmt.Printf(yellowText, "does not exist\n")
		}
	}
	return nil
}
