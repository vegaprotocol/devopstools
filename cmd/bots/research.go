package bots

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/bots"
	"go.uber.org/zap"
)

type ResearchBotsArgs struct {
	*BotsArgs
	BotsAPIToken string
}

var researchBotsArgs ResearchBotsArgs

// researchBotsCmd represents the researchBots command
var researchBotsCmd = &cobra.Command{
	Use:   "research",
	Short: "get research-bots",
	Long:  `get research-bots`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunResearchBots(researchBotsArgs); err != nil {
			researchBotsArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	researchBotsArgs.BotsArgs = &botsArgs

	BotsCmd.AddCommand(researchBotsCmd)
	researchBotsCmd.PersistentFlags().StringVar(&researchBotsArgs.BotsAPIToken, "bots-api-token", "", "API Token for bots endpoint to get private keys")
}

func RunResearchBots(args ResearchBotsArgs) error {
	traders, err := bots.GetResearchBots(args.VegaNetworkName, args.BotsAPIToken)
	if err != nil {
		return err
	}
	byteTraders, err := json.MarshalIndent(traders, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(byteTraders))
	return nil
}
