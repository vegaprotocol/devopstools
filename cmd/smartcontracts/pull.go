package smartcontracts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/types"
	"go.uber.org/zap"
)

type smartContract struct {
	Name       string
	Version    string
	HexAddress string
	Network    types.ETHNetwork
	HomeDir    string
}

type PullArgs struct {
	*SmartContractsArgs
	SmartContracts []smartContract
}

var pullArgs = PullArgs{
	SmartContracts: []smartContract{
		{Name: "MultisigControl", Version: "v1", Network: types.ETHMainnet, HexAddress: "0x164D322B2377C0fdDB73Cd32f24e972A7d9C72F9", HomeDir: "multisigcontrol"},
		{Name: "ERC20Bridge", Version: "v1", Network: types.ETHMainnet, HexAddress: "0xCd403f722b76366f7d609842C589906ca051310f", HomeDir: "erc20bridge"},
		{Name: "ERC20AssetPool", Version: "v1", Network: types.ETHMainnet, HexAddress: "0xF0f0FcDA832415b935802c6dAD0a6dA2c7EAed8f", HomeDir: "erc20assetpool"},
		{Name: "StakingBridge", Version: "v1", Network: types.ETHMainnet, HexAddress: "0x195064D33f09e0c42cF98E665D9506e0dC17de68", HomeDir: "stakingbridge"},
		{Name: "VestingBridge", Version: "v1", Network: types.ETHMainnet, HexAddress: "0x23d1bFE8fA50a167816fBD79D7932577c06011f4", HomeDir: "vestingbridge"},
		{Name: "ERC20BridgeRestricted", Version: "v2", Network: types.ETHMainnet, HexAddress: "0x124Dd8a6044ef048614AEA0AAC86643a8Ae1312D", HomeDir: "erc20bridge"},
		{Name: "MultisigControl", Version: "v2", Network: types.ETHMainnet, HexAddress: "0xDD2df0E7583ff2acfed5e49Df4a424129cA9B58F", HomeDir: "multisigcontrol"},
	},
}

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Get Smart Contracts bytecode and source code from Ethereum Network to local",
	Long:  `Get Smart Contracts bytecode and source code`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunPull(pullArgs); err != nil {
			pullArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	pullArgs.SmartContractsArgs = &smartContractsArgs

	SmartContractsCmd.AddCommand(pullCmd)
}

func RunPull(args PullArgs) error {
	secretStore, err := args.GetServiceSecretStore()
	if err != nil {
		return err
	}
	ethClientManager := ethutils.NewEthereumClientManager(secretStore)
	for _, contract := range args.SmartContracts {
		dir := filepath.Join("smartcontracts", contract.HomeDir, contract.Version)

		if err := ethutils.PullAndStoreSmartContractImmutableData(
			contract.HexAddress, contract.Network, contract.Name, dir, ethClientManager,
		); err != nil {
			return err
		}
		fmt.Printf(" - Downloaded %s(%s) from '%s' and stored in %s\n",
			contract.Name, contract.HexAddress, contract.Network, dir)
	}
	return nil
}
