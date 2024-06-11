package referral

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
	"code.vegaprotocol.io/vega/wallet/wallet"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ProposeArgs struct {
	*Args
	SetupReferralProgram bool
}

var proposeArgs ProposeArgs

var proposeCmd = &cobra.Command{
	Use:   "propose",
	Short: "Propose a referral Program",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunProposeReferral(proposeArgs); err != nil {
			proposeArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	proposeArgs.Args = &args

	Cmd.AddCommand(proposeCmd)
	proposeCmd.PersistentFlags().BoolVar(&proposeArgs.SetupReferralProgram, "setup", false, "Create or update referral program")
}

func RunProposeReferral(args ProposeArgs) error {
	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(ctx, args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Info("Network file loaded", zap.String("name", cfg.EnvironmentName.String()))

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Info("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Info("Looking for healthy gRPC endpoints...")
	healthyEndpoints := tools.FilterHealthyGRPCEndpoints(endpoints)
	if len(healthyEndpoints) == 0 {
		return fmt.Errorf("no healthy gRPC endpoint found on configured data nodes")
	}
	logger.Info("Healthy gRPC endpoints found", zap.Strings("endpoints", healthyEndpoints))

	datanodeClient := datanode.New(healthyEndpoints, 3*time.Second, args.Logger.Named("datanode"))

	logger.Info("Connecting to a datanode's gRPC endpoint...")
	dialCtx, cancelDialing := context.WithTimeout(ctx, 2*time.Second)
	defer cancelDialing()
	datanodeClient.MustDialConnection(dialCtx) // blocking
	logger.Info("Connected to a datanode's gRPC node", zap.String("node", datanodeClient.Target()))

	whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	//
	// Get current referral program
	//
	currentReferralProgram, err := datanodeClient.GetCurrentReferralProgram(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "failed to get current referral program") && strings.Contains(err.Error(), "no rows in result set") {
			logger.Info("Currently there is no referal programm. You can create one.")
			if !args.SetupReferralProgram {
				logger.Warn("You can use --setup to setup referral program")
				return err
			}
		} else {
			return err
		}
	} else {
		logger.Info("Current referral program", zap.Any("config", currentReferralProgram))
	}

	logger.Info("Updating referral program")
	//
	// Setup Referral Program
	//
	if args.SetupReferralProgram {
		err = setupNetworkParametersToSetupReferralProgram(ctx, cfg.EnvironmentName, whaleWallet, datanodeClient, logger)
		if err != nil {
			return err
		}

		networkParameters, err := datanodeClient.ListNetworkParameters(ctx)
		if err != nil {
			return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
		}

		//
		// Prepare Proposal
		//
		minClose, err := time.ParseDuration(networkParameters.Params[netparams.GovernanceProposalReferralProgramMinClose])
		if err != nil {
			return err
		}
		minEnact, err := time.ParseDuration(networkParameters.Params[netparams.GovernanceProposalReferralProgramMinEnact])
		if err != nil {
			return err
		}
		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)
		proposalConfig := programs.NewUpdateReferralProgramProposal(closingTime, enactmentTime)

		//
		// Propose & Vote & Wait
		//
		pubKey := vega.MustFirstKey(whaleWallet)
		err = governance.ProposeVoteAndWait(
			ctx,
			"Referral Program proposal",
			proposalConfig,
			whaleWallet,
			pubKey,
			datanodeClient,
			logger,
		)
		if err != nil {
			return err
		}
	}
	logger.Info("Check API", zap.String("url", fmt.Sprintf("https://api.n00.%s.vega.xyz/api/v2/referral-programs/current", cfg.EnvironmentName)))
	return nil
}

func setupNetworkParametersToSetupReferralProgram(
	ctx context.Context,
	envName config.NetworkName,
	whaleWallet wallet.Wallet,
	datanodeClient *datanode.DataNode,
	logger *zap.Logger,
) error {
	updateParams := map[string]string{
		"governance.proposal.referralProgram.minEnact": "5s",
		"governance.proposal.referralProgram.minClose": "5s",
		"referralProgram.maxReferralTiers":             "10",
		"referralProgram.maxReferralDiscountFactor":    "0.1",
		"referralProgram.maxReferralRewardFactor":      "0.2",
	}
	if envName == config.NetworkDevnet1 {
		updateParams["governance.proposal.referralProgram.requiredParticipation"] = "0.0001"
	}

	networkParameters, err := datanodeClient.ListNetworkParameters(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}

	pubKey := vega.MustFirstKey(whaleWallet)

	updateCount, err := governance.ProposeAndVoteOnNetworkParameters(
		ctx, updateParams, whaleWallet, pubKey, networkParameters, datanodeClient, logger,
	)
	if err != nil {
		return fmt.Errorf("failed to propose and vote for network parameter update proposals: %w", err)
	}

	if updateCount == 0 {
		logger.Info("No network parameter update is required")
		return nil
	}

	return nil
}
