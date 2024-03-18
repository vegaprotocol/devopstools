package multisigcontrol

import (
	"fmt"

	MultisigControl_V1 "github.com/vegaprotocol/devopstools/smartcontracts/multisigcontrol/v1"
	MultisigControl_V2 "github.com/vegaprotocol/devopstools/smartcontracts/multisigcontrol/v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

func DeployMultisigControl(
	version MultisigControlVersion,
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
) (address common.Address, tx *ethTypes.Transaction, err error) {
	switch version {
	case MultisigControlV1:
		address, tx, _, err = MultisigControl_V1.DeployMultisigControl(auth, backend)
	case MultisigControlV2:
		address, tx, _, err = MultisigControl_V2.DeployMultisigControl(auth, backend)
	default:
		err = fmt.Errorf("Invalid Multisig Control Version %s", version)
	}
	return
}
