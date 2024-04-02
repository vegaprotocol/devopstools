package ethutils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetNextNonce(
	ethClient *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
) (uint64, error) {
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return ethClient.PendingNonceAt(context.Background(), address)
}

func TransferEthNoWaitWithNonce(
	ethClient *ethclient.Client,
	fromPrivateKey *ecdsa.PrivateKey,
	toAddress common.Address,
	amount *big.Int,
	nonce uint64,
) (*ethTypes.Transaction, error) {
	errMsg := "failed to send ethereum transaction, %w"

	transactionData, err := GetNextTransactionData(ethClient)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	tx := ethTypes.NewTx(&types.DynamicFeeTx{
		ChainID:   transactionData.chainID,
		Nonce:     nonce,
		GasFeeCap: transactionData.gasFeeCap, // max fee for transaction
		GasTipCap: transactionData.gasTipCap, // max tip
		Gas:       uint64(21000),             // gas limit for a standard ETH transfer is 21000 units
		To:        &toAddress,
		Value:     amount,
		// Data:      txData,
	})

	signedTx, err := ethTypes.SignTx(tx, types.NewLondonSigner(transactionData.chainID), fromPrivateKey)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	if err = ethClient.SendTransaction(context.Background(), signedTx); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return signedTx, nil
}

func TransferEthNoWait(
	ethClient *ethclient.Client,
	fromPrivateKey *ecdsa.PrivateKey,
	toAddress common.Address,
	amount *big.Int,
) (*ethTypes.Transaction, error) {
	nonce, err := GetNextNonce(ethClient, fromPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to send ethereum transaction, %w", err)
	}
	return TransferEthNoWaitWithNonce(ethClient, fromPrivateKey, toAddress, amount, nonce)
}

func TransferEth(
	ethClient *ethclient.Client,
	fromPrivateKey *ecdsa.PrivateKey,
	toAddress common.Address,
	amount *big.Int,
) error {
	tx, err := TransferEthNoWait(ethClient, fromPrivateKey, toAddress, amount)
	if err != nil {
		return err
	}

	receipt, err := bind.WaitMined(context.Background(), ethClient, tx)
	if err != nil {
		return fmt.Errorf("failed to send ethereum transaction, %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("failed to send ethereum transaction, receipt: %+v", receipt)
	}
	return nil
}
