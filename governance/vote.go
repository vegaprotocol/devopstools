package governance

import (
	"context"
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/vegaapi"

	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"
	"code.vegaprotocol.io/vega/wallet/wallet"
)

func VoteOnProposal(ctx context.Context, voteDescription string, proposalId string, voter wallet.Wallet, voterPublicKey string, dataNodeClient vegaapi.DataNodeClient) error {
	request := walletpb.SubmitTransactionRequest{
		Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
			VoteSubmission: &commandspb.VoteSubmission{
				ProposalId: proposalId,
				Value:      vegapb.Vote_VALUE_YES,
			},
		},
	}

	if _, err := walletpkg.SendTransaction(ctx, voter, voterPublicKey, &request, dataNodeClient); err != nil {
		return fmt.Errorf("transaction for proposal %q failed: %w", voteDescription, err)
	}

	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 3)
		vote, err := fetchVoteByProposalIdAndVoter(ctx, proposalId, voterPublicKey, dataNodeClient)
		if err != nil {
			return fmt.Errorf("failed to find vote: %w", err)
		}
		if vote != nil {
			return nil
		}
	}

	return fmt.Errorf("failed to find vote from party %s for proposal %s", voterPublicKey, proposalId)
}

func VoteOnProposalList(ctx context.Context, descriptionToProposalId map[string]string, voter wallet.Wallet, voterPublicKey string, dataNodeClient vegaapi.DataNodeClient) error {
	for description, proposalId := range descriptionToProposalId {
		request := walletpb.SubmitTransactionRequest{
			Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
				VoteSubmission: &commandspb.VoteSubmission{
					ProposalId: proposalId,
					Value:      vegapb.Vote_VALUE_YES,
				},
			},
		}

		if _, err := walletpkg.SendTransaction(ctx, voter, voterPublicKey, &request, dataNodeClient); err != nil {
			return fmt.Errorf("transaction to vote %q failed: %w", description, err)
		}
	}

	proposalIdToVote := map[string]string{}
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 3)
		for _, proposalId := range descriptionToProposalId {
			// skip already fetched proposals
			if _, ok := proposalIdToVote[proposalId]; ok {
				continue
			}
			// fetch proposal by reference
			// on test networks there can be a lot of proposals, so fetching one by one can be more efficient
			vote, err := fetchVoteByProposalIdAndVoter(ctx, proposalId, voterPublicKey, dataNodeClient)
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

func fetchVoteByProposalIdAndVoter(ctx context.Context, proposalId string, voterPartyId string, dataNodeClient vegaapi.DataNodeClient) (*vegapb.Vote, error) {
	voteEdges, err := dataNodeClient.ListVotes(ctx, &v2.ListVotesRequest{
		ProposalId: &proposalId,
		PartyId:    &voterPartyId,
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve votes for proposal %q and party %q: %w", proposalId, voterPartyId, err)
	}

	if len(voteEdges.Votes.Edges) > 1 {
		// This is Vega Network issue
		return nil, fmt.Errorf("found more than 1 vote for proposal %q from same party %q: %+v",
			proposalId, voterPartyId, voteEdges.Votes.Edges,
		)
	} else if len(voteEdges.Votes.Edges) == 1 {
		return voteEdges.Votes.Edges[0].Node, nil
	}
	return nil, nil
}
