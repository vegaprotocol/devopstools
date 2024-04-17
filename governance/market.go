package governance

import (
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/tools"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
)

var LiveMarketStates = []vega.Market_State{
	vega.Market_STATE_ACTIVE,
	vega.Market_STATE_PROPOSED,
	vega.Market_STATE_PENDING,
	vega.Market_STATE_SUSPENDED,
	vega.Market_STATE_SUSPENDED_VIA_GOVERNANCE,
}

func TerminateMarketProposal(closingTime, enactmentTime time.Time, marketName string, marketId string, price string) *commandspb.ProposalSubmission {
	reference := tools.RandAlphaNumericString(40)

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       fmt.Sprintf("Terminate %s market", marketName),
			Description: fmt.Sprintf("Terminate %s market", marketName),
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),

			Change: &vega.ProposalTerms_UpdateMarketState{
				UpdateMarketState: &vega.UpdateMarketState{
					Changes: &vega.UpdateMarketStateConfiguration{
						MarketId:   marketId,
						UpdateType: vega.MarketStateUpdateType_MARKET_STATE_UPDATE_TYPE_TERMINATE,
						Price:      &price,
					},
				},
			},
		},
	}
}
