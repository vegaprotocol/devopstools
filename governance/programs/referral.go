package programs

import (
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

func NewUpdateReferralProgramProposal(closingTime time.Time, enactmentTime time.Time) *commandspb.ProposalSubmission {
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
					Changes: &vega.ReferralProgramChanges{
						BenefitTiers: []*vega.BenefitTier{

							{
								MinimumRunningNotionalTakerVolume: "100000",
								MinimumEpochs:                     "1",
								ReferralRewardFactor:              "0.05",
								ReferralDiscountFactor:            "0.1",
							},
							{
								MinimumRunningNotionalTakerVolume: "1000000",
								MinimumEpochs:                     "1",
								ReferralRewardFactor:              "0.075",
								ReferralDiscountFactor:            "0.1",
							},
							{
								MinimumRunningNotionalTakerVolume: "5000000",
								MinimumEpochs:                     "1",
								ReferralRewardFactor:              "0.1",
								ReferralDiscountFactor:            "0.1",
							},
							{
								MinimumRunningNotionalTakerVolume: "25000000",
								MinimumEpochs:                     "1",
								ReferralRewardFactor:              "0.125",
								ReferralDiscountFactor:            "0.1",
							},
							{
								MinimumRunningNotionalTakerVolume: "75000000",
								MinimumEpochs:                     "1",
								ReferralRewardFactor:              "0.15",
								ReferralDiscountFactor:            "0.1",
							},
							{
								MinimumRunningNotionalTakerVolume: "150000000",
								MinimumEpochs:                     "1",
								ReferralRewardFactor:              "0.175",
								ReferralDiscountFactor:            "0.1",
							},
						},
						StakingTiers: []*vega.StakingTier{
							{
								MinimumStakedTokens:      "100000000000000000000",
								ReferralRewardMultiplier: "1.025",
							},
							{
								MinimumStakedTokens:      "1000000000000000000000",
								ReferralRewardMultiplier: "1.05",
							},
							{
								MinimumStakedTokens:      "5000000000000000000000",
								ReferralRewardMultiplier: "1.1",
							},
							{
								MinimumStakedTokens:      "50000000000000000000000",
								ReferralRewardMultiplier: "1.2",
							},
							{
								MinimumStakedTokens:      "250000000000000000000000",
								ReferralRewardMultiplier: "1.25",
							},
							{
								MinimumStakedTokens:      "500000000000000000000000",
								ReferralRewardMultiplier: "1.3",
							},
						},
						EndOfProgramTimestamp: time.Now().Add(time.Hour * 24 * 365 * 3).Unix(), // TODO: Do We want to have it open almost forever?
						WindowLength:          30,
					},
				},
			},
		},
	}
}
