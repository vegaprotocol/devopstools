package backup

import (
	// "log"

	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ChangeEncryptionKeyArgs struct {
	*BackupRootArgs

	oldPassphrase string
	newPassphrase string

	localStateFile string
}

var changeEncryptionKeyArgs ChangeEncryptionKeyArgs

var changeEncryptionKeyCmd = &cobra.Command{
	Use:   "change-encryption-key",
	Short: "Change an encryption key for the state file",
	Long:  `Takes the state file, decrypt it with the --old-passphrase and encrypts it again with the --new-passphrase.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := DoChangeEncryptionKey(changeEncryptionKeyArgs); err != nil {
			backupArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	changeEncryptionKeyArgs.BackupRootArgs = &backupRootArgs
	changeEncryptionKeyCmd.PersistentFlags().StringVar(&changeEncryptionKeyArgs.oldPassphrase, "old-passphrase", "0123456789abcdef", "The encryption key, the state is encrypted with")
	changeEncryptionKeyCmd.PersistentFlags().StringVar(&changeEncryptionKeyArgs.newPassphrase, "new-passphrase", "0123456789abcdef", "The encryption key, the state should be encryped with")
	changeEncryptionKeyCmd.PersistentFlags().StringVar(&changeEncryptionKeyArgs.localStateFile, "local-state-file", "/tmp/vega-backup-state.json", "Local state file for the vega backup")

	BackupRootCmd.AddCommand(changeEncryptionKeyCmd)
}

func DoChangeEncryptionKey(args ChangeEncryptionKeyArgs) error {
	args.Logger.Info("Reading the state file", zap.String("file", args.localStateFile))

	state, err := LoadFromLocal(args.oldPassphrase, args.localStateFile)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	args.Logger.Info("Updating the encryption key in the state")
	if err := state.UpdateEncryptionKey(args.newPassphrase); err != nil {
		return fmt.Errorf("failed to old passphrase to the new one: %w", err)
	}

	args.Logger.Info("Writing the state file", zap.String("file", args.localStateFile))
	if err := state.WriteLocal(args.localStateFile); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}
