package erc20token

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type InfoArgs struct {
	*ERC20tokenArgs
	VegaNetworkName string
	InfoId          string
}

var infoArgs InfoArgs

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about token",
	Long:  `Get information about token`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunInfo(infoArgs); err != nil {
			infoArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	infoArgs.ERC20tokenArgs = &erc20tokenArgs

	ERC20tokenCmd.AddCommand(infoCmd)
}

func RunInfo(args InfoArgs) error {
	token, err := args.GetToken()
	if err != nil {
		return err
	}
	info, err := token.GetInfo()
	if err != nil {
		return err
	}
	byteInfo, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to parse stats for network '%s', %w", args.VegaNetworkName, err)
	}
	fmt.Println(string(byteInfo))
	return nil
}
