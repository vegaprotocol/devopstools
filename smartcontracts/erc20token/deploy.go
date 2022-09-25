package erc20token

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	ERC20Token_TokenBase "github.com/vegaprotocol/devopstools/smartcontracts/erc20token/TokenBase"
)

func DeployERC20Token(
	version ERC20TokenVersion,
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	name string,
	symbol string,
	decimals uint8,
	totalSupply *big.Int,
) (address common.Address, tx *ethTypes.Transaction, err error) {
	switch version {
	case ERC20TokenBase:
		address, tx, _, err = ERC20Token_TokenBase.DeployTokenBase(auth, backend, name, symbol, decimals, totalSupply)
	default:
		err = fmt.Errorf("Invalid ERC20 Token Version %s", version)
	}
	return
}
