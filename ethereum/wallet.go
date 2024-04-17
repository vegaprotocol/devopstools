package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/vegaprotocol/devopstools/ethutils"

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

	transactOpts, err := ethutils.GetTransactOpts(ctx, ethClient, privateKey)
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
