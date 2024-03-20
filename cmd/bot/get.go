package bot

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vegaprotocol/devopstools/bots"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type GetBotsArgs struct {
	*BotArgs
}

var getBotsArgs GetBotsArgs

// getBotsCmd represents the getBots command
var getBotsCmd = &cobra.Command{
	Use:   "get",
	Short: "get bots",
	Long:  `get bots`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunGetBots(getBotsArgs); err != nil {
			getBotsArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	getBotsArgs.BotArgs = &botArgs

	BotCmd.AddCommand(getBotsCmd)
}

func RunGetBots(args GetBotsArgs) error {
	botsAPIToken := args.BotsAPIToken
	if len(botsAPIToken) == 0 && args.WithSecrets {
		network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
		if err != nil {
			return err
		}
		defer network.Disconnect()

		botsAPIToken = network.BotsApiToken
	}

	traders, err := bots.GetResearchBots(args.VegaNetworkName, botsAPIToken)
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
