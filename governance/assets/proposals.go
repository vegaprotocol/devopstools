package assets

import (
	"math/big"
)

type AssetProposalDetails struct {
	Name     string
	Symbol   string
	Decimals uint64
	Quantum  *big.Int

	ERC20Address             string
	ERC20WithdrawalThreshold string
	ERC20LifetimeLimit       string
}
