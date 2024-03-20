package erc20token

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (t *ERC20Token) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	switch t.Version {
	case ERC20TokenBase:
		return t.v_TokenBase.Mint(opts, to, amount)
	case ERC20TokenOld:
		return t.v_TokenOld.Issue(opts, to, amount)
	case ERC20TokenOther, ERC20TokenMinimal:
		return nil, nil
	default:
		return nil, fmt.Errorf("not implemented for token %s", t.Version)
	}
}
