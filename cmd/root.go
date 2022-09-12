package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	debug  bool
	Logger *zap.Logger
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "devopstools",
	Short: "Scripts to drive Vega Networks",
	Long:  `Manage internal Vega Networks`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		cfg := zap.NewProductionConfig()
		if debug {
			cfg.Level.SetLevel(zap.DebugLevel)
		}
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		Logger, err = cfg.Build()
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
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Print debug logs")
}
