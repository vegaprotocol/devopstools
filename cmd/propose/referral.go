package propose

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/governance/programs"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/veganetwork"

	"code.vegaprotocol.io/vega/core/netparams"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ReferralArgs struct {
	*ProposeArgs
	SetupReferralProgram bool
}

var referralArgs ReferralArgs

// referralCmd represents the referral command
var referralCmd = &cobra.Command{
	Use:   "referral",
	Short: "Referral Program",
	Long:  `Referral Program`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunReferral(referralArgs); err != nil {
			referralArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	referralArgs.ProposeArgs = &proposeArgs

	ProposeCmd.AddCommand(referralCmd)
	referralCmd.PersistentFlags().BoolVar(&referralArgs.SetupReferralProgram, "setup-referral-program", false, "Create or update referral program")
}

func RunReferral(args ReferralArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	var (
		proposerVegawallet = network.VegaTokenWhale
		logger             = args.Logger
	)

	//
	// Get current referral program
	//
	currentReferralProgram, err := network.DataNodeClient.GetCurrentReferralProgram()
	if err != nil {
		if strings.Contains(err.Error(), "failed to get current referral program") && strings.Contains(err.Error(), "no rows in result set") {
			logger.Info("Currently there is no referal programm. You can create one.")
			if !args.SetupReferralProgram {
				logger.Warn("You can use --setup-referral-program to setup referral program")
				return err
			}
		} else {
			return err
		}
	} else {
		logger.Info("Current referral program", zap.Any("config", currentReferralProgram))
	}

	//
	// Setup Referral Program
	//
	if args.SetupReferralProgram {
		err = setupNetworkParametersToSetupReferralProgram(network, logger)
		if err != nil {
			return err
		}
		//
		// Prepare Proposal
		//
		minClose, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalReferralProgramMinClose])
		if err != nil {
			return err
		}
		minEnact, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalReferralProgramMinEnact])
		if err != nil {
			return err
		}
		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)
		proposalConfig := programs.NewUpdateReferralProgramProposal(closingTime, enactmentTime)

		//
		// Propose & Vote & Wait
		//
		err = governance.ProposeVoteAndWait(
			"Referral Program proposal", proposalConfig, proposerVegawallet, network.DataNodeClient, logger,
		)
		if err != nil {
			return err
		}
	}
	logger.Info("Check API", zap.String("url", fmt.Sprintf("https://api.n00.%s.vega.xyz/api/v2/referral-programs/current", network.Network)))
	return nil
}

func setupNetworkParametersToSetupReferralProgram(
	network *veganetwork.VegaNetwork,
	logger *zap.Logger,
) error {
	updateParams := map[string]string{
		"governance.proposal.referralProgram.minEnact": "5s",
		"governance.proposal.referralProgram.minClose": "5s",
		"referralProgram.maxReferralTiers":             "10",
		"referralProgram.maxReferralDiscountFactor":    "0.1",
		"referralProgram.maxReferralRewardFactor":      "0.2",
	}
	if network.Network == types.NetworkDevnet1 {
		updateParams["governance.proposal.referralProgram.requiredParticipation"] = "0.0001"
	}

	updateCount, err := governance.ProposeAndVoteOnNetworkParameters(
		updateParams, network.VegaTokenWhale, network.NetworkParams, network.DataNodeClient, logger,
	)
	if err != nil {
		return err
	}
	if updateCount > 0 {
		if err := network.RefreshNetworkParams(); err != nil {
			return err
		}
	}
	for name, expectedValue := range updateParams {
		if network.NetworkParams.Params[name] != expectedValue {
			return fmt.Errorf("failed to update Network Paramter '%s', current value: '%s', expected value: '%s'",
				name, network.NetworkParams.Params[name], expectedValue,
			)
		}
	}
	return nil
}
