package proposal

import (
	"fmt"
	"time"

	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"
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
	if err := submitTx(proposalDescription, dataNodeClient, proposerVegawallet, logger, &walletTxReq); err != nil {
		return "", err
	}

	//
	// Find Proposal
	//
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 3)
		proposal, err := fetchProposalByReferenceAndProposer(
			reference, proposerVegawallet.PublicKey, dataNodeClient,
		)
		if err != nil {
			return "", err
		}
		if proposal != nil {
			return proposal.Id, nil
		}
	}
	return "", fmt.Errorf("got empty proposal id for the '%s', re %s reference", proposalDescription, reference)
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
		if err := submitTx(description, dataNodeClient, proposerVegawallet, logger, &walletTxReq); err != nil {
			return nil, err
		}
	}

	//
	// Find ProposalIds
	//
	descriptionToProposalId := map[string]string{}
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 3)
		for description, proposalConfig := range descriptionToProposalConfig {
			// skip already fetched proposals
			if _, ok := descriptionToProposalId[description]; ok {
				continue
			}
			// fetch proposal by reference
			// on test networks there can be a lot of proposals, so fetching one by one can be more efficient
			proposal, err := fetchProposalByReferenceAndProposer(
				proposalConfig.Reference, proposerVegawallet.PublicKey, dataNodeClient,
			)
			if err != nil {
				return nil, err
			}
			if proposal != nil {
				descriptionToProposalId[description] = proposal.Id
			}
		}

		if len(descriptionToProposalId) >= len(descriptionToProposalConfig) {
			break
		}
	}

	if len(descriptionToProposalId) < len(descriptionToProposalConfig) {
		return nil, fmt.Errorf("Could not find all proposals, found: \n%+v\n\nall: %+v", descriptionToProposalId, descriptionToProposalConfig)
	}

	return descriptionToProposalId, nil
}

func fetchProposalByReferenceAndProposer(
	reference string,
	proposerPartyId string,
	dataNodeClient vegaapi.DataNodeClient,
) (*vega.Proposal, error) {
	proposalEdges, err := dataNodeClient.ListGovernanceData(&v2.ListGovernanceDataRequest{
		ProposalReference: &reference,
		ProposerPartyId:   &proposerPartyId,
	})
	if err != nil {
		return nil, err
	}
	if len(proposalEdges.Connection.Edges) > 1 {
		// This is Vega Network issue
		return nil, fmt.Errorf("found more than 1 proposal for reference %s: %+v",
			reference, proposalEdges.Connection.Edges,
		)
	} else if len(proposalEdges.Connection.Edges) == 1 {
		proposal := proposalEdges.Connection.Edges[0].Node.Proposal
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
