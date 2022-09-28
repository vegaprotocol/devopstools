package ethutils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//
// TRANSACT OPTS
//

func GetTransactOpts(
	ethClient *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
) (*bind.TransactOpts, error) {
	return GetTransactOptsWithNonce(ethClient, privateKey, 0)
}

func GetTransactOptsWithNonce(
	ethClient *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	nonce uint64,
) (*bind.TransactOpts, error) {
	var (
		errMsg = "failed to get transact opts, %w"
		err    error
	)

	transactionData, err := GetNextTransactionData(ethClient)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	// get transact opts
	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, transactionData.chainID)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	if nonce > 0 {
		transactOpts.Nonce = big.NewInt(int64(nonce))
	}
	// transactOpts.Nonce = big.NewInt(int64(nonce)) // leave empty to use default, i.e. next pending
	// transactOpts.Value = big.NewInt(0) // leave empty to use default 0
	// transactOpts.GasLimit = gasLimit   // leave empty to use estimate
	transactOpts.GasTipCap = transactionData.gasTipCap // Max Priority Fee (tip)
	transactOpts.GasFeeCap = transactionData.gasFeeCap // Max Fee
	return transactOpts, nil
}

//
// TRANSACTION DATA - London Fork
//

type TransactionData struct {
	chainID   *big.Int
	gasFeeCap *big.Int
	gasTipCap *big.Int
}

// Get post London Fork transaction data
func GetNextTransactionData(ethClient *ethclient.Client) (*TransactionData, error) {
	var (
		result = &TransactionData{}
		errMsg = "failed to get next transaction data, %w"
		err    error
	)

	// get chain id
	if result.chainID, err = ethClient.ChainID(context.Background()); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	// get suggested gas price (Base Fee)
	suggestedGasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	// get suggested gas tip cap (Max Priority Fee)
	suggestedGasTipCap, err := ethClient.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	// (!!) TRIPLE (!!) the tip, to speed up transaction
	result.gasTipCap = suggestedGasTipCap.Mul(suggestedGasTipCap, big.NewInt(3))
	// Formula for Max Fee = 2 * Base Fee + Max Priority Fee (tip)
	result.gasFeeCap = new(big.Int).Add(suggestedGasPrice.Mul(suggestedGasPrice, big.NewInt(2)), result.gasTipCap)

	return result, nil
}

func WaitForTransact(
	ethClient *ethclient.Client,
	tx *types.Transaction,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, ethClient, tx)
	if err != nil {
		return fmt.Errorf("transaction failed to mint, %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("Ethereum transaction failed %+v", receipt)
	}
	return nil
}
