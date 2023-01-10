package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Args RootArgs

var RootCmd = &cobra.Command{
	Use:   "devopstools",
	Short: "Scripts to drive Vega Networks",
	Long:  `Manage internal Vega Networks`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		cfg := zap.NewProductionConfig()
		if Args.Debug {
			cfg.Level.SetLevel(zap.DebugLevel)
		}
		// https://github.com/uber-go/zap/issues/584
		cfg.OutputPaths = []string{"stdout"}
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		Args.Logger, err = cfg.Build()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup logger: %s\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run devopstools '%s'\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.PersistentFlags().BoolVar(&Args.Debug, "debug", false, "Print debug logs")
	RootCmd.PersistentFlags().StringVar(&Args.GitHubToken, "github-token", viper.GetString("GITHUB_TOKEN"), "GitHub token to access HashiCorp Vault")
	RootCmd.PersistentFlags().StringVar(&Args.FileWithGitHubToken, "github-token-file", "secret.txt", "file containing GitHub token to access HashiCorp Vault")
	RootCmd.PersistentFlags().StringVar(&Args.HCVaultURL, "hc-vault-url", "https://vault.ops.vega.xyz", "url to HashiCorp Vault")
}
