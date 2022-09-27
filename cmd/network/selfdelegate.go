package network

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"go.uber.org/zap"
)

type SelfDelegateArgs struct {
	*NetworkArgs
}

var selfDelegateArgs SelfDelegateArgs

// selfDelegateCmd represents the selfDelegate command
var selfDelegateCmd = &cobra.Command{
	Use:   "self-delegate",
	Short: "Execute self-delegate for validators",
	Long: `Excecute self-delegate process for every validator, and some steps for non-validators.
	It is safe to call it multiple times, cos it won't repeat steps from previous call.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunSelfDelegate(selfDelegateArgs); err != nil {
			selfDelegateArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	selfDelegateArgs.NetworkArgs = &networkArgs

	NetworkCmd.AddCommand(selfDelegateCmd)
}

func RunSelfDelegate(args SelfDelegateArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	for name, _ := range network.NodeSecrets {
		if err := SelfDelegate(name, network); err != nil {
			return err
		}
	}
	return nil
}

func SelfDelegate(name string, network *veganetwork.VegaNetwork) error {
	minValidatorStake, err := network.NetworkParams.GetMinimumValidatorStake()
	if err != nil {
		return err
	}
	node := network.NodeSecrets[name]

	// General info
	fmt.Printf(" - %s ", name)
	nodeData, isValidator := network.ValidatorsById[node.VegaId]
	if isValidator {
		fmt.Printf("[validator]")
	} else {
		fmt.Printf("[non-validator]")
	}
	fmt.Printf(":\n")

	// ETH
	fmt.Printf("    ethereum balance: ")
	balance, err := network.EthClient.BalanceAt(context.Background(), common.HexToAddress(node.EthereumAddress), nil)
	if err != nil {
		return err
	}
	hBalance := ethutils.WeiToEther(balance)
	fmt.Printf("%f\n", hBalance)

	// STAKE
	fmt.Printf("    stake balance: ")
	balance, err = network.SmartContracts.StakingBridge.GetStakeBalance(node.VegaPubKey)
	if err != nil {
		return err
	}
	hBalance = ethutils.VegaTokenToFullTokens(balance)
	fmt.Printf("%f", hBalance)
	if hBalance.Cmp(minValidatorStake) < 0 {
		diff := new(big.Float)
		diff = diff.Sub(minValidatorStake, hBalance)
		fmt.Printf(" [below required %f]", minValidatorStake)
		if err := Stake(node.VegaPubKey, diff); err != nil {
			return err
		}
	} else {
		fmt.Printf(" [ok]")
	}
	fmt.Printf("\n")

	// SELF-DELEGATE
	if isValidator {
		fmt.Printf("    self-delegate balance: ")

		selfDelegate := new(big.Int)
		var ok bool
		selfDelegate, ok = selfDelegate.SetString(nodeData.StakedByOperator, 0)
		if !ok {
			return fmt.Errorf("failed to convert Staked By Operator '%s' of %s node to big.Int", nodeData.StakedByOperator, name)
		}
		hSelfDelegate := ethutils.VegaTokenToFullTokens(selfDelegate)
		fmt.Printf("%f", hSelfDelegate)
		if isValidator {
			if hSelfDelegate.Cmp(minValidatorStake) < 0 {
				diff := new(big.Float)
				diff = diff.Sub(minValidatorStake, hSelfDelegate)
				fmt.Printf(" [below required %f]", minValidatorStake)
				if err := DelegateToSelf(&node, diff); err != nil {
					return err
				}
			} else {
				fmt.Printf(" [ok]")
			}
		}
		fmt.Printf("\n")
	}

	return nil
}

func Stake(vegaPubKey string, amount *big.Float) error {
	fmt.Printf(" staking %f ...", amount)
	// TODO: implement
	fmt.Printf(" not implemented")
	return nil
}

func DelegateToSelf(node *secrets.VegaNodePrivate, amount *big.Float) error {
	fmt.Printf(" delegating %f ...", amount)
	// TODO: implement
	fmt.Printf(" not implemented")
	return nil
}
