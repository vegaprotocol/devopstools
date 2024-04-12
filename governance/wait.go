package governance

import (
	"context"
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/vegaapi"

	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

func WaitForEnactList(ctx context.Context, descriptionToProposalId map[string]string, dataNodeClient vegaapi.DataNodeClient, logger *zap.Logger) error {
	//
	// Get Latest Enactment
	//
	latestEnactmentTimestamp := int64(0)

	for _, proposalId := range descriptionToProposalId {
		res, err := dataNodeClient.GetGovernanceData(ctx, &v2.GetGovernanceDataRequest{
			ProposalId: &proposalId,
		})
		if err != nil {
			return err
		}

		if slices.Contains(
			[]vega.Proposal_State{vega.Proposal_STATE_FAILED, vega.Proposal_STATE_REJECTED, vega.Proposal_STATE_DECLINED},
			res.Data.Proposal.State,
		) {
			return fmt.Errorf("proposal '%s' is in wrong state %s: %+v", res.Data.Proposal.Rationale.Title, res.Data.Proposal.State.String(), res.Data.Proposal)
		}

		if latestEnactmentTimestamp < res.Data.Proposal.Terms.EnactmentTimestamp {
			latestEnactmentTimestamp = res.Data.Proposal.Terms.EnactmentTimestamp
		}
	}

	//
	// Wait until Latest Enactment
	//
	latestEnactmentTime := time.Unix(latestEnactmentTimestamp, 0)
	untilLatestEnactment := time.Until(latestEnactmentTime)
	if untilLatestEnactment > 5*time.Minute {
		return fmt.Errorf("wait for too long, latest Enactment time is more than 5minutes in the future: %s", latestEnactmentTime)
	} else if untilLatestEnactment > 0 {
		sleepFor := untilLatestEnactment + 5*time.Second
		logger.Info("sleeping unitl latest Enactment time",
			zap.Duration("sleep for", sleepFor),
			zap.Time("sleep until", time.Now().Add(sleepFor)),
		)
		time.Sleep(sleepFor)
	}

	//
	// Validat if every proposal is enacted
	//
	for _, proposalId := range descriptionToProposalId {
		res, err := dataNodeClient.GetGovernanceData(ctx, &v2.GetGovernanceDataRequest{
			ProposalId: &proposalId,
		})
		if err != nil {
			return err
		}

		if res.Data.Proposal.State != vega.Proposal_STATE_ENACTED {
			return fmt.Errorf("proposal '%s' is in wrong state %s: %+v", res.Data.Proposal.Rationale.Title, res.Data.Proposal.State.String(), res.Data.Proposal)
		}
	}

	return nil
}
