package propose

import (
	"fmt"
	"os"
	"strings"
	"time"

	"code.vegaprotocol.io/vega/core/netparams"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/proposal"
	"github.com/vegaprotocol/devopstools/proposal/referral"
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
	Short: "Get Vega Network /statistics",
	Long:  `Get Vega Network /statistics`,
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

	minClose, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinClose])
	if err != nil {
		return err
	}
	minEnact, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinEnact])
	if err != nil {
		return err
	}
	// At this point we know Vega Network is running, and Data Nodes are working

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
		//
		// Propose
		//
		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)
		proposalConfig := referral.NewCreateSimpleReferralSetProposal(closingTime, enactmentTime)

		logger.Info("Submitting Referral Program proposal", zap.Any("proposal", proposalConfig))

		proposalId, err := proposal.SubmitProposal(
			"Setup Referral Program", proposerVegawallet, proposalConfig, network.DataNodeClient, logger,
		)

		if err != nil {
			return fmt.Errorf("failed to propose Referral Program %w", err)
		}

		logger.Info("Proposed Referral Program.", zap.String("proposalId", proposalId))
		//
		// Vote
		//
		logger.Info("Voting on Referral Program", zap.String("proposalId", proposalId))
		err = proposal.VoteOnProposal(
			"Whale vote on Referral Program proposal", proposalId, proposerVegawallet, network.DataNodeClient, logger,
		)
		if err != nil {
			return fmt.Errorf("voting on Referral Program failed %w", err)
		}
		logger.Info("Successfully voted on Referral Program", zap.String("proposalId", proposalId))
	}
	return nil
}
