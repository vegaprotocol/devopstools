package claimcodes

import (
	"fmt"

	ClaimCodes_V1 "github.com/vegaprotocol/devopstools/smartcontracts/claimcodes/v1"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

func DeployClaimCodes(
	version ClaimCodesVersion,
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	vestingBridgeAddress common.Address,
) (address common.Address, tx *ethTypes.Transaction, err error) {
	switch version {
	case ClaimCodesV1:
		address, tx, _, err = ClaimCodes_V1.DeployClaimCodes(auth, backend, vestingBridgeAddress)
	default:
		err = fmt.Errorf("Invalid ERC20 Asset Pool Version %s", version)
	}
	return
}
