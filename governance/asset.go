package governance

import (
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

func NewAssetProposal(closingTime, enactmentTime time.Time) *commandspb.ProposalSubmission {
	reference := tools.RandAlphaNumericString(40)

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "VAP-001 - Create Asset - USD Rewards ($USD-R)",
			Description: `## Summary\n\nProposal to add USD Rewards ($USD-R) as a settlement asset as discussed for incentive iceberg orders\n\n## Rationale\n\n- USD Rewards ($USD-R) will have one market trading on it BTC\n- Given USD Rewards ($USD-R) 18\n- The faucet is open for 1000 every 24 hrs\n- Enjoy the incentive!`,
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:    closingTime.Unix(),
			EnactmentTimestamp:  enactmentTime.Unix(),
			ValidationTimestamp: closingTime.Unix() - 1,
			Change: &vega.ProposalTerms_NewAsset{
				NewAsset: &vega.NewAsset{
					Changes: &vega.AssetDetails{
						Name:     "USD Rewards",
						Symbol:   "USD-R",
						Decimals: 6,
						Quantum:  "1000000",
						Source: &vega.AssetDetails_Erc20{
							Erc20: &vega.ERC20{
								ContractAddress:   "0xaE510B624Ee021a5A80D7cA5bfD8fd039b0B90c8",
								WithdrawThreshold: "1",
								LifetimeLimit:     "7000000000000000000000",
							},
						},
					},
				},
			},
		},
	}
}

// '{
// 	"proposalSubmission": {
// 	 "rationale": {
// 	  "title": "",
// 	  "description": ""
// 	 },
// 	 "terms": {
// 	  "newAsset": {
// 	   "changes": {
// 		"name": "USD Rewards",
// 		"symbol": "USD-R",
// 		"decimals": "6",
// 		"quantum": "1000000",
// 		"erc20": {
// 		 "contractAddress": ,
// 		 "withdrawThreshold": "1",
// 		 "lifetimeLimit": "7000000000000000000000"
// 		}
// 	   }
// 	  },
// 	  "closingTimestamp": 1697231241,
// 	  "enactmentTimestamp": 1697231301,
// 	  "validationTimestamp": 1697231181
// 	 }
// 	}
//    }'
