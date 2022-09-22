package ethutils

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

func ConfirmContractDeployed(
	ctx context.Context,
	client bind.DeployBackend,
	tx *ethTypes.Transaction,
	address common.Address,
) error {
	_, err := bind.WaitDeployed(ctx, client, tx)

	if err != nil {
		return fmt.Errorf("failed to deploy Smart Contract, %w", err)
	}

	byteCode, err := client.CodeAt(ctx, address, nil)
	if err != nil {
		return fmt.Errorf("failed to deploy Smart Contract, cannot get byte code, %w", err)
	}

	if string(byteCode) == "" {
		return fmt.Errorf("failed to deploy Smart Contract, byte code is empty, %w", err)
	}
	return nil
}
