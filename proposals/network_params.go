package proposals

import (
	"fmt"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

func NewUpdateParametersProposal(key string, newValue string, closingTime time.Time, enactmentTime time.Time) *commandspb.ProposalSubmission {
	return &commandspb.ProposalSubmission{
		Reference: tools.RandAlpaNumericString(40),
		Rationale: &vega.ProposalRationale{
			Title:       fmt.Sprintf("Update %s", key),
			Description: fmt.Sprintf("## Summary\n\nChange value of %s to %s from the previous value", key, newValue),
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_UpdateNetworkParameter{
				UpdateNetworkParameter: &vega.UpdateNetworkParameter{
					Changes: &vega.NetworkParameter{
						Key:   key,
						Value: newValue,
					},
				},
			},
		},
	}
}
