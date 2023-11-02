package assets

import (
	"fmt"
	"math/big"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"github.com/vegaprotocol/devopstools/tools"
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

func NewAssetProposal(closingTime, enactmentTime time.Time, details AssetProposalDetails) *commandspb.ProposalSubmission {
	reference := tools.RandAlphaNumericString(40)

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title: fmt.Sprintf("Create Asset - %s ($%s)", details.Name, details.Symbol),
			Description: fmt.Sprintf(`## Summary

Proposal to add USD Rewards ($%s) as a settlement asset as discussed for incentive iceberg orders

## Rationale

- USD Rewards ($%s) will have one market trading on it BTC
- Name: %s
- ERC20 Token(%s)`, details.Symbol, details.Symbol, details.Name, details.ERC20Address),
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:    closingTime.Unix(),
			EnactmentTimestamp:  enactmentTime.Unix(),
			ValidationTimestamp: closingTime.Unix() - 1,
			Change: &vega.ProposalTerms_NewAsset{
				NewAsset: &vega.NewAsset{
					Changes: &vega.AssetDetails{
						Name:     details.Name,
						Symbol:   details.Symbol,
						Decimals: details.Decimals,
						Quantum:  details.Quantum.String(),
						Source: &vega.AssetDetails_Erc20{
							Erc20: &vega.ERC20{
								ContractAddress:   details.ERC20Address,
								WithdrawThreshold: details.ERC20WithdrawalThreshold,
								LifetimeLimit:     details.ERC20LifetimeLimit,
							},
						},
					},
				},
			},
		},
	}
}
