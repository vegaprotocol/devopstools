package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/secrets"
)

type EthWallet struct {
	*secrets.EthereumWalletPrivate

	HexAddress string
	Address    common.Address

	ethClient          *ethclient.Client
	nonce              uint64
	privateKey         *ecdsa.PrivateKey
	cachedTransactOpts *bind.TransactOpts

	txQueue []*types.Transaction
}

func NewEthWallet(
	ethClient *ethclient.Client,
	private *secrets.EthereumWalletPrivate,
) (*EthWallet, error) {
	errMsg := "failed to create new Ethereum Wallet for %s, %w"
	nonce, err := ethClient.PendingNonceAt(context.Background(), common.HexToAddress(private.Address))
	if err != nil {
		return nil, fmt.Errorf(errMsg, private.Address, err)
	}
	privateKey, err := crypto.HexToECDSA(private.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf(errMsg, private.Address, err)
	}
	transactOpts, err := ethutils.GetTransactOpts(ethClient, privateKey)
	if err != nil {
		return nil, fmt.Errorf(errMsg, private.Address, err)
	}
	return &EthWallet{
		EthereumWalletPrivate: private,
		HexAddress:            private.Address,
		Address:               common.HexToAddress(private.Address),
		ethClient:             ethClient,
		nonce:                 nonce,
		privateKey:            privateKey,
		cachedTransactOpts:    transactOpts,
		txQueue:               []*types.Transaction{},
	}, nil
}

func (w *EthWallet) GetNextNonce() uint64 {
	return atomic.AddUint64(&w.nonce, 1) - 1
}

func (w *EthWallet) GetTransactOpts() *bind.TransactOpts {
	var newTransactOptions = *w.cachedTransactOpts
	newTransactOptions.Nonce = big.NewInt(int64(w.GetNextNonce()))

	return &newTransactOptions
}

func (w *EthWallet) ExecuteAndWait(runTransaction func(*bind.TransactOpts) (*types.Transaction, error)) error {
	transactOpts := w.GetTransactOpts()
	tx, err := runTransaction(transactOpts)
	if err != nil {
		return err
	}

	receipt, err := bind.WaitMined(context.Background(), w.ethClient, tx)
	if err != nil {
		return fmt.Errorf("transaction failed to mint, %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("Ethereum transaction failed %+v", receipt)
	}

	return nil
}

func (w *EthWallet) ExecuteAndQueue(runTransaction func(*bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	transactOpts := w.GetTransactOpts()
	tx, err := runTransaction(transactOpts)
	if err != nil {
		return nil, err
	}
	w.txQueue = append(w.txQueue, tx)
	return tx, nil
}

func (w *EthWallet) WaitForQueue() []error {
	txQueue := w.txQueue
	w.txQueue = []*types.Transaction{}

	result := make([]error, len(txQueue))
	for i, tx := range txQueue {
		receipt, err := bind.WaitMined(context.Background(), w.ethClient, tx)
		if err != nil {
			result[i] = fmt.Errorf("transaction failed to mint, %w", err)
		} else if receipt.Status != 1 {
			result[i] = fmt.Errorf("Ethereum transaction failed %+v", receipt)
		} else {
			result[i] = nil
		}
	}

	return result
}
