package party

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/networktools"
	"go.uber.org/zap"
)

type GetStakeArgs struct {
	*PartyArgs
	ValidatorId string
}

var getStakeArgs GetStakeArgs

// getStakeCmd represents the getStake command
var getStakeCmd = &cobra.Command{
	Use:   "get-stake",
	Short: "Get Party Stakes",
	Long:  `Get Party Stakes`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunGetStake(getStakeArgs); err != nil {
			getStakeArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	getStakeArgs.PartyArgs = &partyArgs

	PartyCmd.AddCommand(getStakeCmd)
	getStakeCmd.PersistentFlags().StringVar(&getStakeArgs.ValidatorId, "validator-id", "", "Id of a validator")
}

func RunGetStake(args GetStakeArgs) error {
	network, err := networktools.NewNetworkTools(args.VegaNetworkName, args.Logger)
	if err != nil {
		return err
	}

	api, err := network.GetDataNodeClient()
	if err != nil {
		return err
	}

	stake, err := api.GetPartyTotalStake(args.PartyId)
	if err != nil {
		return err
	}

	fmt.Println(stake)

	return nil
}
