package backup

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/backup"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type PrintConfigArgs struct {
	*BackupArgs
}

var printConfigArgs PrintConfigArgs

var printConfigCmd = &cobra.Command{
	Use:   "print-config",
	Short: "Prints example config",

	Run: func(cmd *cobra.Command, args []string) {
		conf := backup.DefaultConfig()
		confStr, err := conf.Marshal()
		if err != nil {
			printConfigArgs.Logger.Fatal("failed to marshal config", zap.Error(err))
		}

		fmt.Println(confStr)
	},
}

func init() {
	printConfigArgs.BackupArgs = &backupArgs

	BackupCmd.AddCommand(printConfigCmd)
}
