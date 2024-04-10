package governance

import (
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/tools"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
)

// NewBatchProposal either proposals or batchTerms can be nil
// not both at the same time.
func NewBatchProposal(
	title, description string,
	closingTime time.Time,
	// we could use the BatchProposalTermsChange,
	// but instead just do that so we can reuse all existing payloads
	proposals []*commandspb.ProposalSubmission,
	// we also do use the BatchProposalTermsChange just in case
	batchTerms []*vega.BatchProposalTermsChange,
) *commandspb.BatchProposalSubmission {
	var changes []*vega.BatchProposalTermsChange

	for i, v := range proposals {
		if time.Unix(v.Terms.EnactmentTimestamp, 0).Before(closingTime) {
			panic(fmt.Sprintf("proposal at index %v enact(%v) before batching close(%v) time", i, v.Terms.EnactmentTimestamp, closingTime.Unix()))
		}

		change := &vega.BatchProposalTermsChange{
			EnactmentTimestamp: v.Terms.EnactmentTimestamp,
		}

		switch ch := v.Terms.Change.(type) {
		case *vega.ProposalTerms_UpdateMarket:
			change.Change = &vega.BatchProposalTermsChange_UpdateMarket{UpdateMarket: ch.UpdateMarket}
		case *vega.ProposalTerms_NewMarket:
			change.Change = &vega.BatchProposalTermsChange_NewMarket{NewMarket: ch.NewMarket}
		case *vega.ProposalTerms_UpdateNetworkParameter:
			change.Change = &vega.BatchProposalTermsChange_UpdateNetworkParameter{UpdateNetworkParameter: ch.UpdateNetworkParameter}
		case *vega.ProposalTerms_NewFreeform:
			change.Change = &vega.BatchProposalTermsChange_NewFreeform{NewFreeform: ch.NewFreeform}
		case *vega.ProposalTerms_NewSpotMarket:
			change.Change = &vega.BatchProposalTermsChange_NewSpotMarket{NewSpotMarket: ch.NewSpotMarket}
		case *vega.ProposalTerms_UpdateSpotMarket:
			change.Change = &vega.BatchProposalTermsChange_UpdateSpotMarket{UpdateSpotMarket: ch.UpdateSpotMarket}
		case *vega.ProposalTerms_NewTransfer:
			change.Change = &vega.BatchProposalTermsChange_NewTransfer{NewTransfer: ch.NewTransfer}
		case *vega.ProposalTerms_CancelTransfer:
			change.Change = &vega.BatchProposalTermsChange_CancelTransfer{CancelTransfer: ch.CancelTransfer}
		case *vega.ProposalTerms_UpdateMarketState:
			change.Change = &vega.BatchProposalTermsChange_UpdateMarketState{UpdateMarketState: ch.UpdateMarketState}
		case *vega.ProposalTerms_UpdateReferralProgram:
			change.Change = &vega.BatchProposalTermsChange_UpdateReferralProgram{UpdateReferralProgram: ch.UpdateReferralProgram}
		case *vega.ProposalTerms_UpdateVolumeDiscountProgram:
			change.Change = &vega.BatchProposalTermsChange_UpdateVolumeDiscountProgram{UpdateVolumeDiscountProgram: ch.UpdateVolumeDiscountProgram}
		default:
			panic("unsupported market change in batch")
		}

		changes = append(changes, change)
	}

	for i, v := range batchTerms {
		if time.Unix(v.EnactmentTimestamp, 0).Before(closingTime) {
			panic(fmt.Sprintf("batch term at index %v enact(%v) before batching close(%v) time", i, v.EnactmentTimestamp, closingTime.Unix()))
		}

		changes = append(changes, v)
	}

	return &commandspb.BatchProposalSubmission{
		Reference: tools.RandAlphaNumericString(40),
		Rationale: &vega.ProposalRationale{
			Title:       title,
			Description: description,
		},
		Terms: &commandspb.BatchProposalSubmissionTerms{
			ClosingTimestamp: closingTime.Unix(),
			Changes:          changes,
		},
	}
}
