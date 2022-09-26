package ops

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/networktools"
	"go.uber.org/zap"
)

type PsshArgs struct {
	*OpsArgs
	Network           string
	Command           string
	SSHUsername       string
	SSHPrivateKeyfile string
}

var psshArgs PsshArgs

// psshCmd represents the pssh command
var psshCmd = &cobra.Command{
	Use:   "pssh",
	Short: "Run command on every node in network",
	Long:  `Run command on remote machines`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunPSSH(psshArgs); err != nil {
			psshArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	psshArgs.OpsArgs = &opsArgs

	OpsCmd.AddCommand(psshCmd)

	psshCmd.PersistentFlags().StringVar(&psshArgs.Network, "network", "", "Vega Network name")
	if err := psshCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	psshCmd.PersistentFlags().StringVar(&psshArgs.Command, "command", "", "Command to run")
	if err := psshCmd.MarkPersistentFlagRequired("command"); err != nil {
		log.Fatalf("%v\n", err)
	}
	psshCmd.PersistentFlags().StringVar(&psshArgs.SSHUsername, "ssh-user", "", "SSH username")
	if err := psshCmd.MarkPersistentFlagRequired("ssh-user"); err != nil {
		log.Fatalf("%v\n", err)
	}
	psshCmd.PersistentFlags().StringVar(&psshArgs.SSHPrivateKeyfile, "ssh-private-keyfile", "", "File with private SSH key")
	if err := psshCmd.MarkPersistentFlagRequired("ssh-private-keyfile"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunPSSH(args PsshArgs) error {
	network, err := networktools.NewNetworkTools(args.Network, args.Logger)
	if err != nil {
		return err
	}
	sshResults := network.RunCommandOnEveryNode(
		args.SSHUsername, args.SSHPrivateKeyfile, args.Command,
	)

	for host, result := range sshResults {
		if result.Err != nil {
			fmt.Printf("#-ERROR-# %s #-ERROR-#\n%s\n%v\n\n", result.Host, result.Output, result.Err)
		} else {
			fmt.Printf("### %s ###\n%s\n\n", result.Host, result.Output)
		}
		args.Logger.Debug("Execution results", zap.String("host", host), zap.String("result", result.Output), zap.Error(result.Err))
	}
	return nil
}
