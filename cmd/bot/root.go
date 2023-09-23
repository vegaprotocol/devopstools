package bot

import (
	"log"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type BotArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
	BotsAPIToken    string
	WithSecrets     bool
}

var botArgs BotArgs

// Root Command for Bot
var BotCmd = &cobra.Command{
	Use:   "bot",
	Short: "Work with bot",
	Long:  `Work with bot`,
}

func init() {
	botArgs.RootArgs = &rootCmd.Args
	BotCmd.PersistentFlags().StringVar(&botArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := BotCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	BotCmd.PersistentFlags().StringVar(&botArgs.BotsAPIToken, "bots-api-token", "", "API Token for bots endpoint to get private keys")
	BotCmd.PersistentFlags().BoolVar(&botArgs.WithSecrets, "with-secrets", false, "Get the API Token from vault")
}
