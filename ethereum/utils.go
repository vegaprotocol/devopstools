package ethereum

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WaitForTransaction(ctx context.Context, ethClient *ethclient.Client, tx *types.Transaction, timeout time.Duration) error {
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	receipt, err := bind.WaitMined(waitCtx, ethClient, tx)
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed: %+v", receipt)
	}
	return nil
}
