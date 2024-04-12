package propose

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/governance/programs"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/veganetwork"

	"code.vegaprotocol.io/vega/core/netparams"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type VolumeDiscountArgs struct {
	*Args
	SetupVolumeDiscountProgram bool
}

var volumeDiscountArgs VolumeDiscountArgs

var volumeDiscountCmd = &cobra.Command{
	Use:   "volume-discount",
	Short: "Setup Volume Discount Program",
	Long:  `Setup Volume Discount Program`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunVolumeDiscount(volumeDiscountArgs); err != nil {
			volumeDiscountArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	volumeDiscountArgs.Args = &proposeArgs

	Cmd.AddCommand(volumeDiscountCmd)
	volumeDiscountCmd.PersistentFlags().BoolVar(&volumeDiscountArgs.SetupVolumeDiscountProgram, "setup-volume-discount-program", false, "Create or update Volume Discount program")
}

func RunVolumeDiscount(args VolumeDiscountArgs) error {
	ctx := context.Background()
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	proposer := network.VegaTokenWhale
	proposerPublicKey := vega.MustFirstKey(proposer)
	logger := args.Logger

	//
	// Get current volumeDiscount program
	//
	currentVolumeDiscountProgram, err := network.DataNodeClient.GetCurrentVolumeDiscountProgram()
	if err != nil {
		if strings.Contains(err.Error(), "failed to get current volume discount program") && strings.Contains(err.Error(), "no rows in result set") {
			logger.Info("Currently there is no volume discount program. You can create one.")
			if !args.SetupVolumeDiscountProgram {
				logger.Warn("You can use --setup-volume-discount-program to setup Volume Discount program")
				return err
			}
		} else {
			return err
		}
	} else {
		logger.Info("Current volumeDiscount program", zap.Any("config", currentVolumeDiscountProgram))
	}

	//
	// Setup VolumeDiscount Program
	//
	if args.SetupVolumeDiscountProgram {
		err = setupNetworkParametersToSetupVolumeDiscountProgram(ctx, network, logger)
		if err != nil {
			return err
		}
		//
		// Prepare Proposal
		//
		minClose, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalVolumeDiscountProgramMinClose])
		if err != nil {
			return err
		}
		minEnact, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalVolumeDiscountProgramMinEnact])
		if err != nil {
			return err
		}
		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)
		proposalConfig := programs.NewUpdateVolumeDiscountProgramProposal(closingTime, enactmentTime)

		//
		// Propose & Vote & Wait
		//
		err = governance.ProposeVoteAndWait(context.Background(), "Update Volume Discount Program proposal", proposalConfig, proposer, proposerPublicKey, network.DataNodeClient, logger)
		if err != nil {
			return err
		}
	}
	logger.Info("Check API", zap.String("url", fmt.Sprintf("https://api.n00.%s.vega.xyz/api/v2/volume-discount-programs/current", network.Network)))
	return nil
}

func setupNetworkParametersToSetupVolumeDiscountProgram(ctx context.Context, network *veganetwork.VegaNetwork, logger *zap.Logger) error {
	updateParams := map[string]string{
		"governance.proposal.VolumeDiscountProgram.minEnact": "5s",
		"governance.proposal.VolumeDiscountProgram.minClose": "5s",
		"volumeDiscountProgram.maxBenefitTiers":              "10",
		"volumeDiscountProgram.maxVolumeDiscountFactor":      "0.4",
	}

	if network.Network == types.NetworkDevnet1 {
		updateParams["governance.proposal.VolumeDiscountProgram.requiredParticipation"] = "0.0001"
	}

	updateCount, err := governance.ProposeAndVoteOnNetworkParameters(ctx, updateParams, network.VegaTokenWhale, vega.MustFirstKey(network.VegaTokenWhale), network.NetworkParams, network.DataNodeClient, logger)
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
