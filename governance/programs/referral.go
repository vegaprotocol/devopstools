package programs

import (
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

func NewCreateSimpleReferralSetProposal(closingTime time.Time, enactmentTime time.Time) *commandspb.ProposalSubmission {
	return &commandspb.ProposalSubmission{
		Reference: tools.RandAlphaNumericString(40),
		Rationale: &vega.ProposalRationale{
			Title:       "Update the referral program",
			Description: "## Summary\n\nUpdates the referral program",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_UpdateReferralProgram{
				UpdateReferralProgram: &vega.UpdateReferralProgram{
					Changes: &vega.ReferralProgram{
						BenefitTiers: []*vega.BenefitTier{
							{
								MinimumRunningNotionalTakerVolume: "10000",
								MinimumEpochs:                     "6",
								ReferralRewardFactor:              "0.001",
								ReferralDiscountFactor:            "0.001",
							},
							{
								MinimumRunningNotionalTakerVolume: "500000",
								MinimumEpochs:                     "24",
								ReferralRewardFactor:              "0.005",
								ReferralDiscountFactor:            "0.005",
							},
							{
								MinimumRunningNotionalTakerVolume: "1000000",
								MinimumEpochs:                     "48",
								ReferralRewardFactor:              "0.01",
								ReferralDiscountFactor:            "0.01",
							},
						},
						StakingTiers: []*vega.StakingTier{
							{
								MinimumStakedTokens:      "1",
								ReferralRewardMultiplier: "1",
							},
							{
								MinimumStakedTokens:      "2",
								ReferralRewardMultiplier: "2",
							},
							{
								MinimumStakedTokens:      "5",
								ReferralRewardMultiplier: "3",
							},
						},
						EndOfProgramTimestamp: time.Now().Add(time.Hour * 24 * 365 * 3).Unix(), // TODO: Do We want to have it open almost forever?
						WindowLength:          3,
					},
				},
			},
		},
	}
}