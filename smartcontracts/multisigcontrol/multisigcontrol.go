package multisigcontrol

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	MultisigControl_V1 "github.com/vegaprotocol/devopstools/smartcontracts/multisigcontrol/v1"
	MultisigControl_V2 "github.com/vegaprotocol/devopstools/smartcontracts/multisigcontrol/v2"
)

type MultisigControlCommon interface {
	GetCurrentThreshold(opts *bind.CallOpts) (uint16, error)
	GetValidSignerCount(opts *bind.CallOpts) (uint8, error)
	IsNonceUsed(opts *bind.CallOpts, nonce *big.Int) (bool, error)
	IsValidSigner(opts *bind.CallOpts, signer_address common.Address) (bool, error)

	AddSigner(opts *bind.TransactOpts, new_signer common.Address, nonce *big.Int, signatures []byte) (*ethTypes.Transaction, error)
	RemoveSigner(opts *bind.TransactOpts, old_signer common.Address, nonce *big.Int, signatures []byte) (*ethTypes.Transaction, error)
	SetThreshold(opts *bind.TransactOpts, new_threshold uint16, nonce *big.Int, signatures []byte) (*ethTypes.Transaction, error)
	VerifySignatures(opts *bind.TransactOpts, signatures []byte, message []byte, nonce *big.Int) (*ethTypes.Transaction, error)
}

type MultisigControlNewInV2 interface {
	Signers(opts *bind.CallOpts, arg0 common.Address) (bool, error)
	BurnNonce(opts *bind.TransactOpts, nonce *big.Int, signatures []byte) (*ethTypes.Transaction, error)
}

type MultisigControl struct {
	MultisigControlCommon
	MultisigControlNewInV2
	Address common.Address
	Version MultisigControlVersion
	client  *ethclient.Client

	// Minimal implementation
	v1 *MultisigControl_V1.MultisigControl
	v2 *MultisigControl_V2.MultisigControl
}

func NewMultisigControl(
	ethClient *ethclient.Client,
	hexAddress string,
	version MultisigControlVersion,
) (*MultisigControl, error) {
	var err error
	result := &MultisigControl{
		Address: common.HexToAddress(hexAddress),
		Version: version,
		client:  ethClient,
	}
	switch version {
	case MultisigControlV1:
		result.v1, err = MultisigControl_V1.NewMultisigControl(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.MultisigControlCommon = result.v1
	case MultisigControlV2:
		result.v2, err = MultisigControl_V2.NewMultisigControl(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.MultisigControlCommon = result.v2
		result.MultisigControlNewInV2 = result.v2
	}

	return result, nil
}
