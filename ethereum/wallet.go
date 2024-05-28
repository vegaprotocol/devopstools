package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Wallet struct {
	Address common.Address

	nonce              uint64
	privateKey         *ecdsa.PrivateKey
	cachedTransactOpts *bind.TransactOpts
}

func (w *Wallet) GetTransactOpts(ctx context.Context) *bind.TransactOpts {
	nextNonce := atomic.AddUint64(&w.nonce, 1) - 1
	newTransactOptions := *w.cachedTransactOpts
	newTransactOptions.Nonce = big.NewInt(int64(nextNonce))
	newTransactOptions.Context = ctx

	return &newTransactOptions
}

func (w *Wallet) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(data, w.privateKey)
}

func NewWallet(ctx context.Context, ethClient *ethclient.Client, privateHexKey string) (*Wallet, error) {
	privateKey, err := crypto.HexToECDSA(privateHexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key from hex: %w", err)
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)

	nonce, err := ethClient.NonceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account's nonce: %w", err)
	}

	pendingNonce, err := ethClient.PendingNonceAt(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to retrive pending account's nonce: %w", err)
	}

	if nonce != pendingNonce {
		return nil, fmt.Errorf("account's nonce (%d) mismatches account's pending nonce (%d), wallet might already be in use", nonce, pendingNonce)
	}

	transactOpts, err := getTransactOpts(ctx, ethClient, privateKey)
	if err != nil {
		return nil, fmt.Errorf("could not build pre-loaded trasaction options: %w", err)
	}

	return &Wallet{
		Address:            address,
		nonce:              nonce,
		privateKey:         privateKey,
		cachedTransactOpts: transactOpts,
	}, nil
}

func getTransactOpts(ctx context.Context, ethClient *ethclient.Client, privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	return getTransactOptsWithNonce(ctx, ethClient, privateKey, 0)
}

func getTransactOptsWithNonce(ctx context.Context, ethClient *ethclient.Client, privateKey *ecdsa.PrivateKey, nonce uint64) (*bind.TransactOpts, error) {
	var (
		errMsg = "failed to get transact opts, %w"
		err    error
	)

	transactionData, err := getNextTransactionData(ctx, ethClient)
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

type transactionData struct {
	chainID   *big.Int
	gasFeeCap *big.Int
	gasTipCap *big.Int
}

func getNextTransactionData(ctx context.Context, ethClient *ethclient.Client) (*transactionData, error) {
	result := &transactionData{}

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chain ID: %w", err)
	}
	result.chainID = chainID

	// get suggested gas price (Base Fee)
	suggestedGasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested gas price: %w", err)
	}
	// get suggested gas tip cap (Max Priority Fee)
	suggestedGasTipCap, err := ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested gas tip cap: %w", err)
	}

	// (!!) TRIPLE (!!) the tip, to speed up transaction
	result.gasTipCap = suggestedGasTipCap.Mul(suggestedGasTipCap, big.NewInt(3))
	// Formula for Max Fee = 2 * Base Fee + Max Priority Fee (tip)
	result.gasFeeCap = new(big.Int).Add(suggestedGasPrice.Mul(suggestedGasPrice, big.NewInt(2)), result.gasTipCap)

	return result, nil
}
