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
)

func VoteOnProposal(
	voteDescription string,
	proposalId string,
	voterVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	//
	// Vote
	//
	// Prepare vegawallet Transaction Request
	voteWalletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: voterVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
			VoteSubmission: &commandspb.VoteSubmission{
				ProposalId: proposalId,
				Value:      vega.Vote_VALUE_YES,
			},
		},
	}
	if err := submitTx(voteDescription, dataNodeClient, voterVegawallet, logger, &voteWalletTxReq); err != nil {
		return err
	}

	//
	// Find Vote
	//
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 3)
		vote, err := fetchVoteByProposalIdAndVoter(
			proposalId, voterVegawallet.PublicKey, dataNodeClient,
		)
		if err != nil {
			return fmt.Errorf("failed to find vote %w", err)
		}
		if vote != nil {
			return nil
		}
	}
	return fmt.Errorf("failed to find vote from party %s for proposal %s", voterVegawallet.PublicKey, proposalId)
}

func VoteOnProposalList(
	descriptionToProposalId map[string]string,
	voterVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	//
	// Vote
	//
	for description, proposalId := range descriptionToProposalId {
		// Prepare vegawallet Transaction Request
		voteWalletTxReq := walletpb.SubmitTransactionRequest{
			PubKey: voterVegawallet.PublicKey,
			Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
				VoteSubmission: &commandspb.VoteSubmission{
					ProposalId: proposalId,
					Value:      vega.Vote_VALUE_YES,
				},
			},
		}
		if err := submitTx(description, dataNodeClient, voterVegawallet, logger, &voteWalletTxReq); err != nil {
			return err
		}
	}
	//
	// Find Votes
	//
	proposalIdToVote := map[string]string{}
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 3)
		for _, proposalId := range descriptionToProposalId {
			// skip already fetched proposals
			if _, ok := proposalIdToVote[proposalId]; ok {
				continue
			}
			// fetch proposal by reference
			// on test networks there can be a lot of proposals, so fetching one by one can be more efficient
			vote, err := fetchVoteByProposalIdAndVoter(
				proposalId, voterVegawallet.PublicKey, dataNodeClient,
			)
			if err != nil {
				return fmt.Errorf("failed to find vote %w", err)
			}
			if vote != nil {
				proposalIdToVote[proposalId] = vote.PartyId
			}
		}

		if len(proposalIdToVote) >= len(descriptionToProposalId) {
			break
		}
	}
	if len(proposalIdToVote) < len(descriptionToProposalId) {
		return fmt.Errorf("Could not find all proposals, found: \n%+v\n\nall: %+v", proposalIdToVote, descriptionToProposalId)
	}

	return nil
}

func fetchVoteByProposalIdAndVoter(
	proposalId string,
	voterPartyId string,
	dataNodeClient vegaapi.DataNodeClient,
) (*vega.Vote, error) {
	voteEdges, err := dataNodeClient.ListVotes(&v2.ListVotesRequest{
		ProposalId: &proposalId,
		PartyId:    &voterPartyId,
	})
	if err != nil {
		return nil, nil
	}
	if len(voteEdges.Votes.Edges) > 1 {
		// This is Vega Network issue
		return nil, fmt.Errorf("found more than 1 vote for proposalId %s from same party %s: %+v",
			proposalId, voterPartyId, voteEdges.Votes.Edges,
		)
	} else if len(voteEdges.Votes.Edges) == 1 {
		return voteEdges.Votes.Edges[0].Node, nil
	}
	return nil, nil
}
