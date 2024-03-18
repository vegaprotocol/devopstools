package erc20assetpool

import (
	"fmt"

	ERC20AssetPool_V1 "github.com/vegaprotocol/devopstools/smartcontracts/erc20assetpool/v1"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

func DeployERC20AssetPool(
	version ERC20AssetPoolVersion,
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	multisigControlAddress common.Address,
) (address common.Address, tx *ethTypes.Transaction, err error) {
	switch version {
	case ERC20AssetPoolV1:
		address, tx, _, err = ERC20AssetPool_V1.DeployERC20AssetPool(auth, backend, multisigControlAddress)
	default:
		err = fmt.Errorf("Invalid ERC20 Asset Pool Version %s", version)
	}
	return
}
