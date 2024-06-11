package volumediscount

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/governance/programs"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"code.vegaprotocol.io/vega/core/netparams"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ProposeArgs struct {
	*Args
	SetupVolumeDiscountProgram bool
}

var proposeArgs ProposeArgs

var proposeCmd = &cobra.Command{
	Use:   "propose",
	Short: "Setup Volume Discount Program",
	Long:  "Setup Volume Discount Program",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunVolumeDiscount(proposeArgs); err != nil {
			proposeArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	proposeArgs.Args = &args

	Cmd.AddCommand(proposeCmd)
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.SetupVolumeDiscountProgram, "setup", false, "Setup volume discount program. By default, it is dry run.")
}

func RunVolumeDiscount(args ProposeArgs) error {
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
		return fmt.Errorf("no gRPC endpoint found on configured data nodes")
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

	whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	if err := vega.GenerateKeysUpToKey(whaleWallet, cfg.Network.Wallets.VegaTokenWhale.PublicKey); err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	whalePublicKey := cfg.Network.Wallets.VegaTokenWhale.PublicKey

	currentVolumeDiscountProgram, err := datanodeClient.GetCurrentVolumeDiscountProgram(ctx)
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

	if args.SetupVolumeDiscountProgram {
		logger.Debug("Preparing network for volume discount program...")
		updateParams := map[string]string{
			"governance.proposal.VolumeDiscountProgram.minEnact":              "5s",
			"governance.proposal.VolumeDiscountProgram.minClose":              "5s",
			"volumeDiscountProgram.maxBenefitTiers":                           "10",
			"volumeDiscountProgram.maxVolumeDiscountFactor":                   "0.4",
			"governance.proposal.VolumeDiscountProgram.requiredParticipation": "0.0001",
		}
		updatedNetParams, err := vega.UpdateNetworkParameters(ctx, whaleWallet, whalePublicKey, datanodeClient, updateParams, logger)
		if err != nil {
			return fmt.Errorf("failed to prepare network for transfers: %w", err)
		}
		logger.Debug("Network ready for volume discount program")

		//
		// Prepare Proposal
		//
		minClose, err := time.ParseDuration(updatedNetParams.Params[netparams.GovernanceProposalVolumeDiscountProgramMinClose])
		if err != nil {
			return err
		}
		minEnact, err := time.ParseDuration(updatedNetParams.Params[netparams.GovernanceProposalVolumeDiscountProgramMinEnact])
		if err != nil {
			return err
		}
		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)
		proposalConfig := programs.NewUpdateVolumeDiscountProgramProposal(closingTime, enactmentTime)

		//
		// Propose & Vote & Wait
		//
		err = governance.ProposeVoteAndWait(ctx, "Update Volume Discount Program proposal", proposalConfig, whaleWallet, whalePublicKey, datanodeClient, logger)
		if err != nil {
			return err
		}
	}
	logger.Info("Check API", zap.String("url", fmt.Sprintf("https://api.n00.%s.vega.xyz/api/v2/volume-discount-programs/current", cfg.EnvironmentName)))
	return nil
}
