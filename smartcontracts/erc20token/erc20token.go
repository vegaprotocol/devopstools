package erc20token

import (
	"math/big"

	ERC20Token_TokenBase "github.com/vegaprotocol/devopstools/smartcontracts/erc20token/TokenBase"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//
// Functions included in ERC20 standard and implemented by every ERC20 Token
//

type Standard interface {
	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)
	TotalSupply(opts *bind.CallOpts) (*big.Int, error)
	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)
	Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)
	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)
	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)
}

//
// Common functions, but not included in ERC20 standard
//

type Common interface {
	Name(opts *bind.CallOpts) (string, error)
	Symbol(opts *bind.CallOpts) (string, error)
	Decimals(opts *bind.CallOpts) (uint8, error)
}

//
// Base Token extra functionality
//

type Testing interface {
	Faucet(opts *bind.TransactOpts) (*types.Transaction, error)
	Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)
	SetBurnEnabled(opts *bind.TransactOpts, burnEnabled_ bool) (*types.Transaction, error)
	BurnEnabled(opts *bind.CallOpts) (bool, error)
	SetFaucetAmount(opts *bind.TransactOpts, faucetAmount_ *big.Int) (*types.Transaction, error)
	SetFaucetCallLimit(opts *bind.TransactOpts, faucetCallLimit_ *big.Int) (*types.Transaction, error)
	FaucetCallLimit(opts *bind.CallOpts) (*big.Int, error)
	FaucetAmount(opts *bind.CallOpts) (*big.Int, error)

	GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error)
	MINTERROLE(opts *bind.CallOpts) ([32]byte, error)
	GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)
	RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)
	RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)
	HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error)
}

type ERC20Token struct {
	Standard
	Common
	Testing
	Address common.Address

	client *ERC20Token_TokenBase.TokenBase
}

func NewERC20Token(ethClient *ethclient.Client, hexAddress string) (*ERC20Token, error) {
	address := common.HexToAddress(hexAddress)

	client, err := ERC20Token_TokenBase.NewTokenBase(address, ethClient)
	if err != nil {
		return nil, err
	}

	return &ERC20Token{
		Address: address,
		client:  client,
	}, nil
}
