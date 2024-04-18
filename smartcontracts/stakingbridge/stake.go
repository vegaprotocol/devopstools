package stakingbridge

import (
	"math/big"

	"github.com/vegaprotocol/devopstools/tools"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (sb *StakingBridge) StakeBalance(opts *bind.CallOpts, target common.Address, vegaPubKey string) (*big.Int, error) {
	byte32VegaPubKey, err := tools.KeyAsByte32(vegaPubKey)
	if err != nil {
		return nil, err
	}
	return sb.Common.StakeBalance(opts, target, byte32VegaPubKey)
}

func (sb *StakingBridge) Stake(opts *bind.TransactOpts, amount *big.Int, vegaPubKey string) (*types.Transaction, error) {
	byte32VegaPubKey, err := tools.KeyAsByte32(vegaPubKey)
	if err != nil {
		return nil, err
	}
	return sb.Common.Stake(opts, amount, byte32VegaPubKey)
}

func (sb *StakingBridge) RemoveStake(opts *bind.TransactOpts, amount *big.Int, vegaPubKey string) (*types.Transaction, error) {
	byte32VegaPubKey, err := tools.KeyAsByte32(vegaPubKey)
	if err != nil {
		return nil, err
	}
	return sb.Common.RemoveStake(opts, amount, byte32VegaPubKey)
}
