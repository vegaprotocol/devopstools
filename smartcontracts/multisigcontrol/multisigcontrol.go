package multisigcontrol

import (
	"context"
	"fmt"
	"math/big"

	MultisigControl_V1 "github.com/vegaprotocol/devopstools/smartcontracts/multisigcontrol/v1"
	MultisigControl_V2 "github.com/vegaprotocol/devopstools/smartcontracts/multisigcontrol/v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Common interface {
	GetCurrentThreshold(opts *bind.CallOpts) (uint16, error)
	GetValidSignerCount(opts *bind.CallOpts) (uint8, error)
	IsNonceUsed(opts *bind.CallOpts, nonce *big.Int) (bool, error)
	IsValidSigner(opts *bind.CallOpts, signer_address common.Address) (bool, error)

	AddSigner(opts *bind.TransactOpts, new_signer common.Address, nonce *big.Int, signatures []byte) (*ethTypes.Transaction, error)
	RemoveSigner(opts *bind.TransactOpts, old_signer common.Address, nonce *big.Int, signatures []byte) (*ethTypes.Transaction, error)
	SetThreshold(opts *bind.TransactOpts, new_threshold uint16, nonce *big.Int, signatures []byte) (*ethTypes.Transaction, error)
	VerifySignatures(opts *bind.TransactOpts, signatures []byte, message []byte, nonce *big.Int) (*ethTypes.Transaction, error)
}

type NewInV2 interface {
	Signers(opts *bind.CallOpts, arg0 common.Address) (bool, error)
	BurnNonce(opts *bind.TransactOpts, nonce *big.Int, signatures []byte) (*ethTypes.Transaction, error)
}

type MultisigControl struct {
	Common
	NewInV2
	Address common.Address
	Version Version
	client  *ethclient.Client

	// Minimal implementation
	v1 *MultisigControl_V1.MultisigControl
	v2 *MultisigControl_V2.MultisigControl
}

func (m *MultisigControl) GetSigners(ctx context.Context) ([]common.Address, error) {
	latestBlockNumber, err := m.client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %w", err)
	}

	signerCounter := map[common.Address]int{}

	// Increase counter with every addition
	switch m.Version {
	case V1:
		addedIterator, err := m.v1.FilterSignerAdded(&bind.FilterOpts{
			Start:   0,
			End:     &latestBlockNumber,
			Context: ctx,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to filter SignerAdded events: %w", err)
		}
		for addedIterator.Next() {
			signerCounter[addedIterator.Event.NewSigner] += 1
		}
	case V2:
		addedIterator, err := m.v2.FilterSignerAdded(&bind.FilterOpts{
			Start:   0,
			End:     &latestBlockNumber,
			Context: ctx,
		})
		if err != nil {
			return nil, fmt.Errorf(" failed to filter SignerAdded events: %w", err)
		}
		for addedIterator.Next() {
			signerCounter[addedIterator.Event.NewSigner] += 1
		}
	default:
		return nil, fmt.Errorf("version %q is not supported", m.Version)
	}

	// Decrease counter with every removal
	switch m.Version {
	case V1:
		removedIterator, err := m.v1.FilterSignerRemoved(&bind.FilterOpts{
			Start:   0,
			End:     &latestBlockNumber,
			Context: ctx,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to filter SignerRemoved events: %w", err)
		}
		for removedIterator.Next() {
			signerCounter[removedIterator.Event.OldSigner] -= 1
		}
	case V2:
		removedIterator, err := m.v2.FilterSignerRemoved(&bind.FilterOpts{
			Start:   0,
			End:     &latestBlockNumber,
			Context: ctx,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to filter SignerRemoved events: %w", err)
		}
		for removedIterator.Next() {
			signerCounter[removedIterator.Event.OldSigner] -= 1
		}
	default:
		return nil, fmt.Errorf("version '%s' is not supported", m.Version)
	}

	result := []common.Address{}

	for signerAddress, counter := range signerCounter {
		if counter == 1 {
			result = append(result, signerAddress)
		} else if counter != 0 {
			return nil, fmt.Errorf(
				"failed to get signers, counter for '%s' signer is '%d'; it should be 0 or 1",
				signerAddress, counter,
			)
		}
	}

	return result, nil
}

func NewMultisigControl(
	ethClient *ethclient.Client,
	hexAddress string,
	version Version,
) (*MultisigControl, error) {
	var err error
	result := &MultisigControl{
		Address: common.HexToAddress(hexAddress),
		Version: version,
		client:  ethClient,
	}
	switch version {
	case V1:
		result.v1, err = MultisigControl_V1.NewMultisigControl(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.Common = result.v1
	case V2:
		result.v2, err = MultisigControl_V2.NewMultisigControl(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.Common = result.v2
		result.NewInV2 = result.v2
	}

	return result, nil
}
