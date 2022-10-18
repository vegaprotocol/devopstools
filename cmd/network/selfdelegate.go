package network

import (
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/wallet"
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
	minValidatorStake, err := network.NetworkParams.GetMinimumValidatorStake()
	if err != nil {
		return err
	}
	humanMinValidatorStake := ethutils.VegaTokenToFullTokens(minValidatorStake)
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
		logger.Info("staking to node", zap.String("node", name), zap.String("vegaPubKey", node.VegaPubKey),
			zap.String("amount", humanStakeAmount.String()), zap.String("stakingBridgeAddress", stakingBridge.Address.Hex()))
		tx, err := stakingBridge.Stake(opts, stakeAmount, node.VegaPubKey)
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

	//
	// Delegate
	//
	lastBlockData, err := network.DataNodeClient.LastBlockData()
	if err != nil {
		return err
	}
	statistics, err := network.DataNodeClient.Statistics()
	if err != nil {
		return err
	}

	resultsChannel := make(chan error, len(network.ValidatorsById))
	var wg sync.WaitGroup
	//for id, validator := range network.ValidatorsById {
	for name, nodeSecrets := range network.NodeSecrets {
		validator, isValidator := network.ValidatorsById[nodeSecrets.VegaId]
		if !isValidator {
			continue
		}
		stakedByOperator, ok := new(big.Int).SetString(validator.StakedByOperator, 0)
		if !ok {
			logger.Error("failed to parse StakedByOperator", zap.String("node", name),
				zap.String("StakedByOperator", validator.StakedByOperator), zap.Error(err))
		}
		if minValidatorStake.Cmp(stakedByOperator) <= 0 {
			logger.Info("node already self-delegated enough", zap.String("node", name),
				zap.String("StakedByOperator", validator.StakedByOperator))
			continue
		}
		partyTotalStake, err := network.DataNodeClient.GetPartyTotalStake(nodeSecrets.VegaId)
		if err != nil {
			return err
		}
		if partyTotalStake.Cmp(minValidatorStake) < 0 {
			logger.Warn("party doesn't have visible stake yet - you might need to wait till next epoch", zap.String("node", name),
				zap.String("partyTotalStake", partyTotalStake.String()))
			// TODO: write wait functionality
			//continue
		}

		wg.Add(1)
		go func(name string, validator *vega.Node, nodeSecrets secrets.VegaNodePrivate, lastBlockData *vegaapipb.LastBlockHeightResponse, chainID string) {
			defer wg.Done()
			vegawallet, err := wallet.NewVegaWallet(&secrets.VegaWalletPrivate{
				Id:             nodeSecrets.VegaId,
				PublicKey:      nodeSecrets.VegaPubKey,
				PrivateKey:     nodeSecrets.VegaPrivateKey,
				RecoveryPhrase: nodeSecrets.VegaRecoveryPhrase,
			})
			if err != nil {
				logger.Error("failed to create wallet", zap.String("node", name), zap.Error(err))
				resultsChannel <- fmt.Errorf("failed to create wallet for %s node", name)
				return
			}
			walletTxReq := walletpb.SubmitTransactionRequest{
				PubKey: nodeSecrets.VegaPubKey,
				Command: &walletpb.SubmitTransactionRequest_DelegateSubmission{
					DelegateSubmission: &commandspb.DelegateSubmission{
						NodeId: nodeSecrets.VegaId,
						Amount: minValidatorStake.String(),
					},
				},
			}

			signedTx, err := vegawallet.SignTxWithPoW(&walletTxReq, lastBlockData)
			if err != nil {
				logger.Error("failed to sign a trasnaction", zap.String("node", name), zap.Error(err))
				resultsChannel <- fmt.Errorf("failed to sign a transaction for %s node", name)
				return
			}

			submitReq := &vegaapipb.SubmitTransactionRequest{
				Tx:   signedTx,
				Type: vegaapipb.SubmitTransactionRequest_TYPE_SYNC,
			}
			submitResponse, err := network.DataNodeClient.SubmitTransaction(submitReq)
			if err != nil {
				logger.Error("failed to submit a trasnaction", zap.String("node", name), zap.Error(err))
				resultsChannel <- fmt.Errorf("failed to submit a transaction for %s node", name)
				return
			}
			if !submitResponse.Success {
				logger.Error("transaction submission failure", zap.String("node", name), zap.Error(err))
				resultsChannel <- fmt.Errorf("transaction submission failure for %s node", name)
				return
			}
			logger.Info("successful delegation", zap.String("node", name), zap.Any("response", submitResponse))
		}(name, validator, nodeSecrets, lastBlockData, statistics.Statistics.ChainId)
	}
	wg.Wait()
	close(resultsChannel)

	var failureCount int

	for err := range resultsChannel {
		if err != nil {
			failureCount += 1
		}
	}
	successCount := len(network.ValidatorsById) - failureCount
	logger.Info("Delegate summary", zap.Int("successCount", successCount), zap.Int("failureCount", failureCount))
	if failureCount > 0 {
		return fmt.Errorf("failed to self-delegate to all of the parties")
	}
	fmt.Printf("DONE\n")

	return nil
}
