package remote

import (
	"path"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/tools"
)

type RemoteArgs struct {
	*rootCmd.RootArgs

	ServerHost string
	ServerUser string
	ServerKey  string
}

var remoteArgs RemoteArgs

var RemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Set of commands to perform actions on a remote servers",
}

func init() {
	remoteArgs.RootArgs = &rootCmd.Args

	RemoteCmd.PersistentFlags().StringVar(&remoteArgs.ServerHost, "host", "", "Host to connect to")
	RemoteCmd.PersistentFlags().StringVar(&remoteArgs.ServerUser, "user", tools.CurrentUsername(), "User to connect to the server")
	RemoteCmd.PersistentFlags().StringVar(&remoteArgs.ServerKey, "key-path", path.Join(tools.UserHomeDir(), ".ssh", "id_rsa"), "Key to authenticate given user on the server")

	if err := RemoteCmd.MarkPersistentFlagRequired("host"); err != nil {
		panic(err)
	}
}
