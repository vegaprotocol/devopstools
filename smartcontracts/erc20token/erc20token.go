package erc20token

import (
	"math/big"

	ERC20Token_IERC20 "github.com/vegaprotocol/devopstools/smartcontracts/erc20token/IERC20"
	ERC20Token_TokenBase "github.com/vegaprotocol/devopstools/smartcontracts/erc20token/TokenBase"
	ERC20Token_TokenOld "github.com/vegaprotocol/devopstools/smartcontracts/erc20token/TokenOld"
	ERC20Token_TokenOther "github.com/vegaprotocol/devopstools/smartcontracts/erc20token/TokenOther"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//
// Functions included in ERC20 standard and implemented by every ERC20 Token
//

type ERC20TokenStandard interface {
	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)
	TotalSupply(opts *bind.CallOpts) (*big.Int, error)
	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)
	Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)
	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)
}

//
// Common functions, but not included in ERC20 standard
//

type ERC20TokenCommon interface {
	Name(opts *bind.CallOpts) (string, error)
	Symbol(opts *bind.CallOpts) (string, error)
	Decimals(opts *bind.CallOpts) (uint8, error)
}

//
// Base Token extra functionality
//

type ERC20TokenTesting interface {
	Faucet(opts *bind.TransactOpts) (*types.Transaction, error)
	Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)
	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)
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
	ERC20TokenStandard
	ERC20TokenCommon
	ERC20TokenTesting
	Address common.Address
	Version ERC20TokenVersion
	client  *ethclient.Client

	// Minimal implementation
	v_IERC20 *ERC20Token_IERC20.IERC20
	// Most common implementation
	v_TokenOther *ERC20Token_TokenOther.TokenOther
	// deprecated - don't ever use
	v_TokenOld *ERC20Token_TokenOld.TokenOld
	// For our testing implementation
	v_TokenBase *ERC20Token_TokenBase.TokenBase
}

func NewERC20Token(ethClient *ethclient.Client, hexAddress string, version ERC20TokenVersion) (*ERC20Token, error) {
	var err error
	result := &ERC20Token{
		Address: common.HexToAddress(hexAddress),
		Version: version,
		client:  ethClient,
	}
	switch version {
	case ERC20TokenMinimal:
		result.v_IERC20, err = ERC20Token_IERC20.NewIERC20(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.ERC20TokenStandard = result.v_IERC20
	case ERC20TokenOther:
		result.v_TokenOther, err = ERC20Token_TokenOther.NewTokenOther(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.ERC20TokenStandard = result.v_TokenOther
		result.ERC20TokenCommon = result.v_TokenOther
	case ERC20TokenOld:
		result.v_TokenOld, err = ERC20Token_TokenOld.NewTokenOld(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.ERC20TokenStandard = result.v_TokenOld
		result.ERC20TokenCommon = result.v_TokenOld
	case ERC20TokenBase:
		result.v_TokenBase, err = ERC20Token_TokenBase.NewTokenBase(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.ERC20TokenStandard = result.v_TokenBase
		result.ERC20TokenCommon = result.v_TokenBase
		result.ERC20TokenTesting = result.v_TokenBase
	}

	return result, nil
}
