package network

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
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
	logger := args.Logger
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	//
	// STAKE
	//
	var (
		vegaToken              = network.SmartContracts.VegaToken
		stakingBridge          = network.SmartContracts.StakingBridge
		minterWallet           = network.NetworkMainWallet
		humanMissingStakeTotal = new(big.Float)
		humanMissingStake      = map[string]*big.Float{}
		status                 string

		stakeSuccessCount, stakeFailureCount, stakeOKCount int
	)
	humanMinValidatorStake, err := network.NetworkParams.GetMinimumValidatorStake()
	if err != nil {
		return err
	}
	vegaTokenInfo, err := vegaToken.GetInfo()
	if err != nil {
		return fmt.Errorf("failed to get vega token info %s: %w", vegaToken.Address, err)
	}
	// Gather missing stake
	for name, node := range network.NodeSecrets {
		stakeBalance, err := stakingBridge.GetStakeBalance(node.VegaPubKey)
		if err != nil {
			return err
		}
		humanStakeBalance := ethutils.VegaTokenToFullTokens(stakeBalance)
		if humanStakeBalance.Cmp(humanMinValidatorStake) < 0 {
			humanMissingStake[name] = new(big.Float).Sub(humanMinValidatorStake, humanStakeBalance)
			humanMissingStakeTotal = humanMissingStakeTotal.Add(humanMissingStakeTotal, humanMissingStake[name])
			status = "stake below minimum"
		} else {
			stakeOKCount += 1
			status = "stake balance ok"
		}
		logger.Info(fmt.Sprintf("Smart Contract state: %s", status), zap.String("node", name),
			zap.String("stakeBalance", humanStakeBalance.String()),
			zap.String("requiredBalance", humanMinValidatorStake.String()),
			zap.String("vegaPubKey", node.VegaPubKey))
	}

	//
	// Prepare MinterWallet: Mint enough tokens and Increase Allowance
	//
	var (
		mintTx      *ethTypes.Transaction
		allowanceTx *ethTypes.Transaction
	)
	minterTokenBalance, err := vegaToken.BalanceOf(&bind.CallOpts{}, minterWallet.Address)
	if err != nil {
		return fmt.Errorf("failed to get balance, %w", err)
	}
	requiredTokensTotal := ethutils.VegaTokenFromFullTokens(humanMissingStakeTotal)
	if minterTokenBalance.Cmp(requiredTokensTotal) < 0 {
		diff := new(big.Int).Sub(requiredTokensTotal, minterTokenBalance)
		opts := minterWallet.GetTransactOpts()
		logger.Info("minting Vega Token", zap.String("amount", diff.String()), zap.String("token", vegaTokenInfo.Name),
			zap.String("ethWallet", minterWallet.Address.Hex()),
			zap.String("balanceBefore", minterTokenBalance.String()),
			zap.String("increasingBy", diff.String()),
			zap.String("requiredBalance", requiredTokensTotal.String()),
			zap.String("tokenAddress", vegaToken.Address.Hex()))
		mintTx, err = vegaToken.Mint(opts, minterWallet.Address, diff)
		if err != nil {
			return fmt.Errorf("failed to mint %s of %s (%s) for %s: %w", diff.String(),
				vegaTokenInfo.Name, vegaToken.Address.Hex(), minterWallet.Address.Hex(), err)
		}
	} else {
		logger.Info("no need to mint", zap.String("token", vegaTokenInfo.Name),
			zap.String("ethWallet", minterWallet.Address.Hex()),
			zap.String("currentBalance", minterTokenBalance.String()),
			zap.String("requiredBalance", requiredTokensTotal.String()),
			zap.String("tokenAddress", vegaToken.Address.Hex()))
	}
	minterTokenAllowance, err := vegaToken.Allowance(&bind.CallOpts{}, minterWallet.Address, stakingBridge.Address)
	if err != nil {
		return fmt.Errorf("failed to get allowance for staking bridge %s, %w", stakingBridge.Address.Hex(), err)
	}
	if minterTokenAllowance.Cmp(requiredTokensTotal) < 0 {
		diff := new(big.Int).Sub(requiredTokensTotal, minterTokenAllowance)
		opts := minterWallet.GetTransactOpts()
		logger.Info("increasing allowance", zap.String("amount", diff.String()), zap.String("token", vegaTokenInfo.Name),
			zap.String("ethWallet", minterWallet.Address.Hex()),
			zap.String("allowanceBefore", minterTokenAllowance.String()),
			zap.String("increasingBy", diff.String()),
			zap.String("requiredAllowance", requiredTokensTotal.String()),
			zap.String("tokenAddress", vegaToken.Address.Hex()))
		allowanceTx, err = vegaToken.IncreaseAllowance(opts, stakingBridge.Address, diff)
		if err != nil {
			return fmt.Errorf("failed to increase allowance: %w", err)
		}
	} else {
		logger.Info("no need to increase allowance", zap.String("token", vegaTokenInfo.Name),
			zap.String("ethWallet", minterWallet.Address.Hex()),
			zap.String("currentAllowance", minterTokenAllowance.String()),
			zap.String("requiredAllowance", requiredTokensTotal.String()))
	}
	// wait
	if mintTx != nil {
		if err = ethutils.WaitForTransaction(network.EthClient, mintTx, time.Minute); err != nil {
			logger.Error("failed to mint", zap.String("token", vegaTokenInfo.Name), zap.Error(err))
			return fmt.Errorf("transaction failed to mints: %w", err)
		}
		logger.Info("successfully minted", zap.String("ethWallet", minterWallet.Address.Hex()),
			zap.String("token", vegaTokenInfo.Name), zap.String("tokenAddress", vegaToken.Address.Hex()))
	}
	if allowanceTx != nil {
		if err = ethutils.WaitForTransaction(network.EthClient, allowanceTx, time.Minute); err != nil {
			logger.Error("failed to increase allowance", zap.String("ethWallet", minterWallet.Address.Hex()),
				zap.String("token", vegaTokenInfo.Name), zap.String("tokenAddress", vegaToken.Address.Hex()), zap.Error(err))
			return fmt.Errorf("transaction failed to increase allowance: %w", err)
		}
		logger.Info("successfully increased allowance", zap.String("ethWallet", minterWallet.Address.Hex()),
			zap.String("token", vegaTokenInfo.Name), zap.String("tokenAddress", vegaToken.Address.Hex()))
	}

	//
	// Stake to nodes
	//
	var (
		stakeTxs = map[string]*ethTypes.Transaction{}
	)
	for name, humanStakeAmount := range humanMissingStake {
		var (
			node        = network.NodeSecrets[name]
			stakeAmount = ethutils.VegaTokenFromFullTokens(humanStakeAmount)
			opts        = minterWallet.GetTransactOpts()
		)
		bytePubKey, err := hex.DecodeString(node.VegaPubKey)
		if err != nil {
			return err
		}
		var byte32PubKey [32]byte
		copy(byte32PubKey[:], bytePubKey)
		logger.Info("staking to node", zap.String("node", name), zap.String("vegaPubKey", node.VegaPubKey),
			zap.String("amount", humanStakeAmount.String()), zap.String("stakingBridgeAddress", stakingBridge.Address.Hex()))
		tx, err := stakingBridge.Stake(opts, stakeAmount, byte32PubKey)
		if err != nil {
			stakeFailureCount += 1
			logger.Error("failed to stake", zap.String("node", name), zap.Error(err))
		} else {
			stakeTxs[name] = tx
		}
	}
	// wait
	for name, tx := range stakeTxs {
		if tx == nil {
			continue
		}
		logger.Debug("waiting", zap.Any("tx", tx))
		if err = ethutils.WaitForTransaction(network.EthClient, tx, time.Minute); err != nil {
			stakeFailureCount += 1
			logger.Error("failed to stake", zap.String("node", name),
				zap.Any("tx", tx), zap.Error(err))
		} else {
			stakeSuccessCount += 1
			logger.Info("successfully staked", zap.String("node", name))
		}
	}
	logger.Info("Stake summary", zap.Int("success", stakeSuccessCount), zap.Int("failure", stakeFailureCount),
		zap.Int("already ok", stakeOKCount))

	// for name, _ := range network.NodeSecrets {
	// 	if err := SelfDelegate(name, network); err != nil {
	// 		return err
	// 	}
	// }
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
