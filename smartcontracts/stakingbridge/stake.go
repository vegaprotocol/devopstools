package stakingbridge

import (
	"context"
	"fmt"
	"math/big"

	"github.com/vegaprotocol/devopstools/ethutils"
	StakingBridge_V1 "github.com/vegaprotocol/devopstools/smartcontracts/stakingbridge/v1"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (sb *StakingBridge) GetStakeBalance(
	vegaPubKey string,
) (balance *big.Int, err error) {
	var (
		latestBlockNumber uint64
		depositedIterator *StakingBridge_V1.StakingBridgeStakeDepositedIterator
		removedIterator   *StakingBridge_V1.StakingBridgeStakeRemovedIterator
	)
	byte32VegaPubKey, err := ethutils.VegaPubKeyToByte32(vegaPubKey)
	if err != nil {
		return nil, err
	}
	// get last block
	latestBlockNumber, err = sb.client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get balance for '%s', failed to get latest block number, %w", vegaPubKey, err)
	}

	balance = big.NewInt(0)

	// Increase balance with every deposit
	depositedIterator, err = sb.v1.FilterStakeDeposited(&bind.FilterOpts{
		Start:   0,
		End:     &latestBlockNumber,
		Context: context.Background(),
	}, nil, [][32]byte{byte32VegaPubKey})
	if err != nil {
		return nil, fmt.Errorf("failed to get balance for '%s', failed to Filter Stake Deposited events, %w", vegaPubKey, err)
	}
	for depositedIterator.Next() {
		balance = balance.Add(balance, depositedIterator.Event.Amount)
	}

	// Decrease balance with every remove
	removedIterator, err = sb.v1.FilterStakeRemoved(&bind.FilterOpts{
		Start:   0,
		End:     &latestBlockNumber,
		Context: context.Background(),
	}, nil, [][32]byte{byte32VegaPubKey})
	if err != nil {
		return nil, fmt.Errorf("failed to get balance for '%s', failed to Filter Stake Removed events, %w", vegaPubKey, err)
	}
	for removedIterator.Next() {
		// fmt.Printf("remove %#v\n", depositedIterator.Event)
		balance = balance.Sub(balance, removedIterator.Event.Amount)
	}

	// Ignore Stake Transfer events - transfer does not change Stake Balance of vegawallet

	return
}

func (sb *StakingBridge) StakeBalance(opts *bind.CallOpts, target common.Address, vegaPubKey string) (*big.Int, error) {
	byte32VegaPubKey, err := ethutils.VegaPubKeyToByte32(vegaPubKey)
	if err != nil {
		return nil, err
	}
	return sb.Common.StakeBalance(opts, target, byte32VegaPubKey)
}

func (sb *StakingBridge) Stake(opts *bind.TransactOpts, amount *big.Int, vegaPubKey string) (*types.Transaction, error) {
	byte32VegaPubKey, err := ethutils.VegaPubKeyToByte32(vegaPubKey)
	if err != nil {
		return nil, err
	}
	return sb.Common.Stake(opts, amount, byte32VegaPubKey)
}

func (sb *StakingBridge) RemoveStake(opts *bind.TransactOpts, amount *big.Int, vegaPubKey string) (*types.Transaction, error) {
	byte32VegaPubKey, err := ethutils.VegaPubKeyToByte32(vegaPubKey)
	if err != nil {
		return nil, err
	}
	return sb.Common.RemoveStake(opts, amount, byte32VegaPubKey)
}
