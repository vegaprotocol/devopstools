package erc20bridge

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	ERC20Bridge_V1 "github.com/vegaprotocol/devopstools/smartcontracts/erc20bridge/v1"
	ERC20Bridge_V2 "github.com/vegaprotocol/devopstools/smartcontracts/erc20bridge/v2"
)

func DeployERC20Bridge(
	version ERC20BridgeVersion,
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	multisigControlAddress common.Address,
	erc20AssetPoolAddress common.Address,
) (address common.Address, tx *ethTypes.Transaction, err error) {
	switch version {
	case "v1":
		address, tx, _, err = ERC20Bridge_V1.DeployERC20Bridge(auth, backend, multisigControlAddress, erc20AssetPoolAddress)
	case "v2":
		address, tx, _, err = ERC20Bridge_V2.DeployERC20BridgeRestricted(auth, backend, erc20AssetPoolAddress)
	default:
		err = fmt.Errorf("Invalid ERC20 Bridge Version %s", version)
	}
	return
}
