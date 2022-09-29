package claimcodes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	ClaimCodes_V1 "github.com/vegaprotocol/devopstools/smartcontracts/claimcodes/v1"
)

type ClaimCodesCommon interface {
	AllowedCountries(opts *bind.CallOpts, arg0 [2]byte) (bool, error)
	Commitments(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error)
	Controller(opts *bind.CallOpts) (common.Address, error)
	Issuers(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error)
	AllowCountries(opts *bind.TransactOpts, countries [][2]byte) (*types.Transaction, error)
	BlockCountries(opts *bind.TransactOpts, countries [][2]byte) (*types.Transaction, error)
	ClaimTargeted(opts *bind.TransactOpts, sig ClaimCodes_V1.Signature, clm ClaimCodes_V1.Claim, country [2]byte, target common.Address) (*types.Transaction, error)
	ClaimUntargeted(opts *bind.TransactOpts, sig ClaimCodes_V1.Signature, clm ClaimCodes_V1.Claim, country [2]byte) (*types.Transaction, error)
	CommitUntargeted(opts *bind.TransactOpts, s [32]byte) (*types.Transaction, error)
}

type ClaimCodes struct {
	ClaimCodesCommon
	Address common.Address
	Version ClaimCodesVersion
	client  *ethclient.Client

	// Minimal implementation
	v1 *ClaimCodes_V1.ClaimCodes
}

func NewClaimCodes(
	ethClient *ethclient.Client,
	hexAddress string,
	version ClaimCodesVersion,
) (*ClaimCodes, error) {
	var err error
	result := &ClaimCodes{
		Address: common.HexToAddress(hexAddress),
		Version: version,
		client:  ethClient,
	}
	switch version {
	case ClaimCodesV1:
		result.v1, err = ClaimCodes_V1.NewClaimCodes(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.ClaimCodesCommon = result.v1
	}

	return result, nil
}
