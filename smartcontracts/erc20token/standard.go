package erc20token

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (t *ERC20Token) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	return t.client.BalanceOf(opts, account)
}

func (t *ERC20Token) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	return t.client.TotalSupply(opts)
}

func (t *ERC20Token) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.client.Approve(opts, spender, amount)
}

func (t *ERC20Token) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.client.Transfer(opts, to, amount)
}

func (t *ERC20Token) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.client.TransferFrom(opts, from, to, amount)
}

func (t *ERC20Token) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	return t.client.Allowance(opts, owner, spender)
}

func (t *ERC20Token) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return t.client.IncreaseAllowance(opts, spender, addedValue)
}
