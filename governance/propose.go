package governance

import (
	"context"
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/governance/networkparameters"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vegaapi"

	"code.vegaprotocol.io/vega/core/netparams"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"
	"code.vegaprotocol.io/vega/wallet/wallet"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

const IgnoreProposer = ""

func SubmitProposalList(ctx context.Context, descriptionToProposalConfig map[string]*commandspb.ProposalSubmission, proposer wallet.Wallet, proposerPublicKey string, dataNodeClient vegaapi.DataNodeClient) (map[string]string, error) {
	for description, proposalConfig := range descriptionToProposalConfig {
		request := walletpb.SubmitTransactionRequest{
			Command: &walletpb.SubmitTransactionRequest_ProposalSubmission{
				ProposalSubmission: proposalConfig,
			},
		}
		if _, err := walletpkg.SendTransaction(ctx, proposer, proposerPublicKey, &request, dataNodeClient); err != nil {
			return nil, fmt.Errorf("transaction for proposal %q failed: %w", description, err)
		}
	}

	//
	// Find ProposalIds
	//
	descriptionToProposalId := map[string]string{}

	for description, proposalConfig := range descriptionToProposalConfig {
		// skip already fetched proposals
		if _, ok := descriptionToProposalId[description]; ok {
			continue
		}
		// fetch proposal by reference
		// on test networks there can be a lot of proposals, so fetching one by one can be more efficient

		proposalId, err := tools.RetryReturn(6, 10*time.Second, func() (string, error) {
			proposal, err := fetchProposalByReferenceAndProposer(ctx, proposalConfig.Reference, dataNodeClient)
			if err != nil {
				return "", fmt.Errorf("failed to find proposal: %w", err)
			}
			if proposal != nil {
				return proposal.Id, nil
			}

			return "", fmt.Errorf("got empty proposal id for the '%s', re %s reference", description, proposalConfig.Reference)
		})
		if err != nil {
			return nil, fmt.Errorf("failed to find proposal: %w", err)
		}

		descriptionToProposalId[description] = proposalId
	}

	if len(descriptionToProposalId) < len(descriptionToProposalConfig) {
		return nil, fmt.Errorf("Could not find all proposals, found: \n%+v\n\nall: %+v", descriptionToProposalId, descriptionToProposalConfig)
	}

	return descriptionToProposalId, nil
}

func fetchProposalByReferenceAndProposer(ctx context.Context, reference string, dataNodeClient vegaapi.DataNodeClient) (*vegapb.Proposal, error) {
	res, err := tools.RetryReturn(6, 10*time.Second, func() (*v2.GetGovernanceDataResponse, error) {
		return dataNodeClient.GetGovernanceData(ctx, &v2.GetGovernanceDataRequest{
			Reference: &reference,
		})
	})
	if err != nil {
		return nil, err
	}
	if res != nil {
		proposal := res.Data.Proposal
		if slices.Contains(
			[]vegapb.Proposal_State{vegapb.Proposal_STATE_FAILED, vegapb.Proposal_STATE_REJECTED, vegapb.Proposal_STATE_DECLINED},
			proposal.State,
		) {
			return nil, fmt.Errorf("proposal '%s' is in wrong state %s: %+v", proposal.Rationale.Title, proposal.State.String(), proposal)
		}
		return proposal, nil
	}
	return nil, nil
}

func ProposeAndVoteOnNetworkParameters(ctx context.Context, desiredValues map[string]string, proposer wallet.Wallet, proposerPublicKey string, networkParams *types.NetworkParams, dataNodeClient vegaapi.DataNodeClient, logger *zap.Logger) (int64, error) {
	minClose, err := time.ParseDuration(networkParams.Params[netparams.GovernanceProposalUpdateNetParamMinClose])
	if err != nil {
		return 0, fmt.Errorf("could not parse network parameter %q", netparams.GovernanceProposalUpdateNetParamMinClose)
	}
	minEnact, err := time.ParseDuration(networkParams.Params[netparams.GovernanceProposalUpdateNetParamMinEnact])
	if err != nil {
		return 0, fmt.Errorf("could not parse network parameter %q", netparams.GovernanceProposalUpdateNetParamMinEnact)
	}

	closingTime := time.Now().Add(time.Second * 20).Add(minClose)
	enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)

	descriptionToProposalConfig := map[string]*commandspb.ProposalSubmission{}
	for name, desiredValue := range desiredValues {
		currentValue, ok := networkParams.Params[name]
		if ok && desiredValue == currentValue {
			logger.Debug("Not including network parameter to proposal as already set to the updated value",
				zap.String("name", name),
				zap.String("value", currentValue),
			)
			continue
		}

		logger.Debug("Including network parameter to proposal to the updated value",
			zap.String("name", name),
			zap.String("current-value", currentValue),
			zap.String("desired-value", desiredValue),
			zap.Time("proposal-closing-time", closingTime),
			zap.Time("proposal-enactment-time", enactmentTime),
		)

		description := fmt.Sprintf("Update Network Paramter %q=%q", name, desiredValue)
		descriptionToProposalConfig[description] = networkparameters.NewUpdateParametersProposal(name, desiredValue, closingTime, enactmentTime)
	}

	if len(descriptionToProposalConfig) == 0 {
		logger.Debug("No network parameter to update")
		return 0, nil
	}

	if err := ProposeVoteAndWaitList(ctx, descriptionToProposalConfig, proposer, proposerPublicKey, dataNodeClient, logger); err != nil {
		return 0, fmt.Errorf("proposals failed: %w", err)
	}

	return int64(len(descriptionToProposalConfig)), nil
}

func FindProposalID(ctx context.Context, proposerPubKey string, reference string, client vegaapi.DataNodeClient) (string, error) {
	res, err := client.GetGovernanceData(ctx, &v2.GetGovernanceDataRequest{
		Reference: &reference,
	})
	if err != nil {
		return "", fmt.Errorf("failed to find any proposal with the %s reference: %w", reference, err)
	}

	if res != nil {
		proposal := res.Data.Proposal
		if slices.Contains(
			[]vegapb.Proposal_State{vegapb.Proposal_STATE_FAILED, vegapb.Proposal_STATE_REJECTED, vegapb.Proposal_STATE_DECLINED},
			proposal.State,
		) {
			return "", fmt.Errorf("proposal '%s' is in wrong state %s: %+v", proposal.Rationale.Title, proposal.State.String(), proposal)
		}

		// If We provided any proposer check if it matches
		if proposerPubKey != IgnoreProposer && proposerPubKey != proposal.PartyId {
			return "", fmt.Errorf(
				"found proposal for reference(%s) but proposer does not match: expected %s, got %s",
				reference,
				proposerPubKey,
				proposal.PartyId,
			)
		}

		return proposal.Id, nil
	}

	return "", fmt.Errorf("failed to find proposal for reference %s", reference)
}
