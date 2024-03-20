package governance

import (
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/governance/networkparameters"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"github.com/vegaprotocol/devopstools/wallet"

	"code.vegaprotocol.io/vega/core/netparams"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

func SubmitProposal(
	proposalDescription string,
	proposerVegawallet *wallet.VegaWallet,
	proposal *commandspb.ProposalSubmission,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) (string, error) {
	reference := proposal.Reference
	//
	// Propose
	//
	// Prepare vegawallet Transaction Request
	walletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_ProposalSubmission{
			ProposalSubmission: proposal,
		},
	}
	if err := SubmitTx(proposalDescription, dataNodeClient, proposerVegawallet, logger, &walletTxReq); err != nil {
		return "", err
	}

	//
	// Find Proposal
	//
	proposalId, err := tools.RetryReturn(6, 10*time.Second, func() (string, error) {
		proposal, err := fetchProposalByReferenceAndProposer(reference, dataNodeClient)
		if err != nil {
			return "", fmt.Errorf("failed to find proposal: %w", err)
		}
		if proposal != nil {
			return proposal.Id, nil
		}

		return "", fmt.Errorf("got empty proposal id for the '%s', re %s reference", proposalDescription, reference)
	})

	return proposalId, fmt.Errorf("failed to find proposal: %w", err)
}

func SubmitProposalList(
	descriptionToProposalConfig map[string]*commandspb.ProposalSubmission,
	proposerVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) (map[string]string, error) {
	//
	// Propose
	//
	for description, proposalConfig := range descriptionToProposalConfig {
		// Prepare vegawallet Transaction Request
		walletTxReq := walletpb.SubmitTransactionRequest{
			PubKey: proposerVegawallet.PublicKey,
			Command: &walletpb.SubmitTransactionRequest_ProposalSubmission{
				ProposalSubmission: proposalConfig,
			},
		}
		if err := SubmitTx(description, dataNodeClient, proposerVegawallet, logger, &walletTxReq); err != nil {
			return nil, err
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
			proposal, err := fetchProposalByReferenceAndProposer(proposalConfig.Reference, dataNodeClient)
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

func fetchProposalByReferenceAndProposer(reference string, dataNodeClient vegaapi.DataNodeClient) (*vega.Proposal, error) {
	res, err := tools.RetryReturn(6, 10*time.Second, func() (*v2.GetGovernanceDataResponse, error) {
		return dataNodeClient.GetGovernanceData(&v2.GetGovernanceDataRequest{
			Reference: &reference,
		})
	})
	if err != nil {
		return nil, err
	}
	if res != nil {
		proposal := res.Data.Proposal
		if slices.Contains(
			[]vega.Proposal_State{vega.Proposal_STATE_FAILED, vega.Proposal_STATE_REJECTED, vega.Proposal_STATE_DECLINED},
			proposal.State,
		) {
			return nil, fmt.Errorf("proposal '%s' is in wrong state %s: %+v", proposal.Rationale.Title, proposal.State.String(), proposal)
		}
		return proposal, nil
	}
	return nil, nil
}

func ProposeAndVoteOnNetworkParameters(
	nameToValue map[string]string,
	proposerVegawallet *wallet.VegaWallet,
	networkParams *types.NetworkParams,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) (int64, error) {
	errorMsgPrefix := "failed to Propose and Vote on Update Network Paramter"
	minClose, err := time.ParseDuration(networkParams.Params[netparams.GovernanceProposalUpdateNetParamMinClose])
	if err != nil {
		return 0, fmt.Errorf("%s, %w", errorMsgPrefix, err)
	}
	minEnact, err := time.ParseDuration(networkParams.Params[netparams.GovernanceProposalUpdateNetParamMinEnact])
	if err != nil {
		return 0, fmt.Errorf("%s, %w", errorMsgPrefix, err)
	}
	logger.Debug("proposal limits", zap.Duration("minClose", minClose), zap.Duration("minEnact", minEnact))

	//
	// Prepare proposals
	//
	descriptionToProposalConfig := map[string]*commandspb.ProposalSubmission{}
	for name, value := range nameToValue {
		if currentValue, ok := networkParams.Params[name]; ok && value == currentValue {
			logger.Info("Skip Networ Paramter proposal", zap.String("name", name), zap.String("value", value))
			continue
		}
		description := fmt.Sprintf("Update Network Paramter '%s'='%s'", name, value)
		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)
		descriptionToProposalConfig[description] = networkparameters.NewUpdateParametersProposal(
			name, value, closingTime, enactmentTime,
		)
	}
	if len(descriptionToProposalConfig) == 0 {
		return 0, nil
	}

	//
	// Propose & Vote
	//
	err = ProposeVoteAndWaitList(
		descriptionToProposalConfig, proposerVegawallet, dataNodeClient, logger,
	)
	if err != nil {
		return 0, fmt.Errorf("%s, %w", errorMsgPrefix, err)
	}

	return int64(len(descriptionToProposalConfig)), nil
}

func WaitForNetworkParameters(network *veganetwork.VegaNetwork, expectedNetworkParameters map[string]string) error {
	return tools.RetryRun(6, 10*time.Second, func() error {
		if err := network.RefreshNetworkParams(); err != nil {
			return fmt.Errorf("failed to refresh network parameters: %w", err)
		}

		for key, expectedValue := range expectedNetworkParameters {
			currentValue, ok := network.NetworkParams.Params[key]
			if !ok {
				return fmt.Errorf("the %s network parameter not found yet on the network", key)
			}

			if currentValue != expectedValue {
				return fmt.Errorf(
					"the %s network parameter not updated yet: expected value %s, current value: %s",
					key,
					expectedValue,
					currentValue,
				)
			}
		}

		return nil
	})
}
