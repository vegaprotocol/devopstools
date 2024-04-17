package bot

import (
	"log"

	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
	VegaNetworkName string
	BotsAPIToken    string
	WithSecrets     bool
}

var botArgs Args

var Cmd = &cobra.Command{
	Use:   "bot",
	Short: "Work with bot",
	Long:  `Work with bot`,
}

func init() {
	botArgs.RootArgs = &rootCmd.Args
	Cmd.PersistentFlags().StringVar(&botArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := Cmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	Cmd.PersistentFlags().StringVar(&botArgs.BotsAPIToken, "bots-api-token", "", "API Token for bots endpoint to get private keys")
	Cmd.PersistentFlags().BoolVar(&botArgs.WithSecrets, "with-secrets", false, "Get the API Token from vault")
}
