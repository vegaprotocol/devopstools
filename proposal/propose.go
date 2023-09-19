package proposal

import (
	"fmt"
	"time"

	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
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
	// PROPOSE
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
	time.Sleep(time.Second * 10)
	res, err := dataNodeClient.ListGovernanceData(&v2.ListGovernanceDataRequest{
		ProposalReference: &reference,
	})
	if err != nil {
		return "", err
	}
	var proposalId string
	for _, edge := range res.Connection.Edges {
		if edge.Node.Proposal.Reference == reference {
			proposalId = edge.Node.Proposal.Id
			break
		}
	}

	if len(proposalId) < 1 {
		return "", fmt.Errorf("got empty proposal id for the '%s', re %s reference", proposalDescription, reference)
	}
	return proposalId, nil
}
