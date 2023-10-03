package programs

import (
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

func NewUpdateVolumeDiscountProgramProposal(closingTime time.Time, enactmentTime time.Time) *commandspb.ProposalSubmission {
	return &commandspb.ProposalSubmission{
		Reference: tools.RandAlphaNumericString(40),
		Rationale: &vega.ProposalRationale{
			Title:       "Update the Volume Discount program",
			Description: "## Summary\n\nUpdates the volume discount program",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_UpdateVolumeDiscountProgram{
				UpdateVolumeDiscountProgram: &vega.UpdateVolumeDiscountProgram{
					Changes: &vega.VolumeDiscountProgram{
						BenefitTiers: []*vega.VolumeBenefitTier{
							{
								MinimumRunningNotionalTakerVolume: "1000",
								VolumeDiscountFactor:              "0.001",
							},
							{
								MinimumRunningNotionalTakerVolume: "20000",
								VolumeDiscountFactor:              "0.002",
							},
							{
								MinimumRunningNotionalTakerVolume: "30000",
								VolumeDiscountFactor:              "0.003",
							},
						},
						EndOfProgramTimestamp: time.Now().Add(time.Hour * 24 * 365 * 3).Unix(), // TODO: Do We want to have it open almost forever?
						WindowLength:          7,
					},
				},
			},
		},
	}
}
