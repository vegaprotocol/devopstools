package propose

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/governance/networkparameters"

	"code.vegaprotocol.io/vega/core/netparams"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type NetworkParameterArgs struct {
	*ProposeArgs
	NetworkParameterName  string
	NetworkParameterValue string
}

var networkParameterArgs NetworkParameterArgs

// networkParameterCmd represents the networkParameter command
var networkParameterCmd = &cobra.Command{
	Use:   "network-parameter",
	Short: "Get or Update (propose & vote) on Network Paramter",
	Long:  `Get or Update (propose & vote) on Network Paramter≠`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunNetworkParameter(networkParameterArgs); err != nil {
			networkParameterArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	networkParameterArgs.ProposeArgs = &proposeArgs

	ProposeCmd.AddCommand(networkParameterCmd)
	networkParameterCmd.PersistentFlags().StringVar(&networkParameterArgs.NetworkParameterName, "name", "", "Network Paramter name")
	if err := networkParameterCmd.MarkPersistentFlagRequired("name"); err != nil {
		log.Fatalf("%v\n", err)
	}
	networkParameterCmd.PersistentFlags().StringVar(&networkParameterArgs.NetworkParameterValue, "set-value", "", "Leave empty to get current value. Set to update parameter to a new value")
}

func RunNetworkParameter(args NetworkParameterArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	var (
		proposerVegawallet = network.VegaTokenWhale
		logger             = args.Logger
	)
	minClose, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalUpdateNetParamMinClose])
	if err != nil {
		return err
	}
	minEnact, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalUpdateNetParamMinEnact])
	if err != nil {
		return err
	}

	if len(args.NetworkParameterValue) == 0 {
		res, err := network.DataNodeClient.ListNetworkParameters(&v2.ListNetworkParametersRequest{})
		if err != nil {
			return err
		}
		var value string = ""
		for _, dd := range res.NetworkParameters.Edges {
			if dd.Node.Key == args.NetworkParameterName {
				value = dd.Node.Value
				break
			}
		}
		if len(value) > 0 {
			fmt.Printf("Network '%s' paramter value: '%s'\n", args.NetworkParameterName, value)
		} else {
			fmt.Printf("There is no network partamter with name '%s'\n", args.NetworkParameterName)
		}
	} else {
		//
		// Propose
		//
		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)
		proposalConfig := networkparameters.NewUpdateParametersProposal(
			args.NetworkParameterName, args.NetworkParameterValue, closingTime, enactmentTime,
		)

		logger.Info("Submitting Update Network Paramter proposal", zap.Any("proposal", proposalConfig))

		proposalId, err := governance.SubmitProposal(
			"Update Network Parameter", proposerVegawallet, proposalConfig, network.DataNodeClient, logger,
		)
		if err != nil {
			return fmt.Errorf("failed to propose Update Network Paramter %w", err)
		}

		logger.Info("Proposed Update Network Paramter.", zap.String("proposalId", proposalId))
		//
		// Vote
		//
		logger.Info("Voting on Update Network Paramter", zap.String("proposalId", proposalId))
		err = governance.VoteOnProposal(
			"Whale vote on Update Network Paramter proposal", proposalId, proposerVegawallet, network.DataNodeClient, logger,
		)
		if err != nil {
			return fmt.Errorf("voting on Update Network Paramter failed %w", err)
		}
		logger.Info("Successfully voted on Update Network Paramter", zap.String("proposalId", proposalId))
	}

	return nil
}
