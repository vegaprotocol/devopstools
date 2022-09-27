package stakingbridge

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	StakingBridge_V1 "github.com/vegaprotocol/devopstools/smartcontracts/stakingbridge/v1"
)

func (sb *StakingBridge) GetStakeBalance(
	vegaPubKey string,
) (balance *big.Int, err error) {
	var (
		byte32VegaPubKey  [32]byte
		latestBlockNumber uint64
		depositedIterator *StakingBridge_V1.StakingBridgeStakeDepositedIterator
		removedIterator   *StakingBridge_V1.StakingBridgeStakeRemovedIterator
	)
	byteVegaPubKey, err := hex.DecodeString(vegaPubKey)
	if err != nil {
		return nil, err
	}
	copy(byte32VegaPubKey[:], byteVegaPubKey)
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
		//fmt.Printf("remove %#v\n", depositedIterator.Event)
		balance = balance.Sub(balance, removedIterator.Event.Amount)
	}

	// Ignore Stake Transfer events - transfer does not change Stake Balance of vegawallet

	return
}
