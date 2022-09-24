package erc20token

import (
	"fmt"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/smartcontracts"
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20token"
	"github.com/vegaprotocol/devopstools/types"
)

type ERC20tokenArgs struct {
	*rootCmd.RootArgs
	Address      string
	EthNetwork   string
	TokenVersion string
	TokenName    string
	EthereumURL  string
}

var erc20tokenArgs ERC20tokenArgs

// Root Command for OPS
var ERC20tokenCmd = &cobra.Command{
	Use:   "erc20token",
	Short: "General erc20token tasks",
	Long:  `General erc20token tasks`,
}

func init() {
	erc20tokenArgs.RootArgs = &rootCmd.Args

	ERC20tokenCmd.PersistentFlags().StringVar(&erc20tokenArgs.TokenName, "token", "", "Name of a token")
	ERC20tokenCmd.PersistentFlags().StringVar(&erc20tokenArgs.Address, "address", "", "Token address")
	ERC20tokenCmd.PersistentFlags().StringVar(&erc20tokenArgs.EthNetwork, "eth-network", "sepolia", "Used with address, specify which Ethereum Network to use")
	ERC20tokenCmd.PersistentFlags().StringVar(&erc20tokenArgs.TokenVersion, "token-version", "TokenBase", "Used with address, specify version of a token, allowed values: TokenBase, TokenOther and TokenMinimal")

	ERC20tokenCmd.PersistentFlags().StringVar(&erc20tokenArgs.EthereumURL, "ethereum-url", "", "Optional URL to connect to ethereum network, e.g. infura")
}

func (ra *ERC20tokenArgs) GetSmartContractsManager() (*smartcontracts.SmartContractsManager, error) {
	if len(ra.EthereumURL) > 0 {
		return smartcontracts.NewSmartContractsManagerWithEthURL(ra.EthereumURL), nil
	} else {
		return ra.RootArgs.GetSmartContractsManager()
	}
}

func (ra *ERC20tokenArgs) GetToken() (*erc20token.ERC20Token, error) {
	smartContractManager, err := ra.GetSmartContractsManager()
	if err != nil {
		return nil, err
	}
	if len(ra.TokenName) > 0 {
		return smartContractManager.GetAssetWithName(ra.TokenName)
	}
	if len(ra.Address) > 0 {
		if len(ra.EthNetwork) == 0 {
			return nil, fmt.Errorf("need to provide --eth-network when using --address")
		}
		ethNetwork := types.ETHNetwork(ra.EthNetwork)
		if err = ethNetwork.IsValid(); err != nil {
			return nil, fmt.Errorf("wrong eth network %s, %w", ra.EthNetwork, err)
		}
		tokenVersion := erc20token.ERC20TokenVersion(ra.TokenVersion)
		if err = tokenVersion.IsValid(); err != nil {
			return nil, fmt.Errorf("wrong token version %s, %w", ra.TokenVersion, err)
		}
		return smartContractManager.GetAsset(
			ra.Address, ethNetwork, tokenVersion,
		)
	}
	return nil, fmt.Errorf("need to provide token or address")
}
