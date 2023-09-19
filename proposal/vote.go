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
	// VOTE
	//
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
	time.Sleep(time.Second * 10)
	res, err := dataNodeClient.ListVotes(&v2.ListVotesRequest{
		ProposalId: &proposalId,
		PartyId:    &voterVegawallet.PublicKey,
	})
	if err != nil {
		return fmt.Errorf("failed to find vote %w", err)
	}
	if len(res.Votes.Edges) < 1 {
		return fmt.Errorf("failed to find vote from party %s for proposal %s", voterVegawallet.PublicKey, proposalId)
	}
	return nil
}
