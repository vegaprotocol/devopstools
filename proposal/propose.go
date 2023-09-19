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

	var prop *vega.Proposal
	for _, edge := range res.Connection.Edges {
		if edge.Node.Proposal.Reference == reference {
			prop = edge.Node.Proposal
			break
		}
	}

	if prop == nil {
		return "", fmt.Errorf("got empty proposal id for the '%s', re %s reference", proposalDescription, reference)
	}
	if slices.Contains(
		[]vega.Proposal_State{vega.Proposal_STATE_FAILED, vega.Proposal_STATE_REJECTED, vega.Proposal_STATE_DECLINED},
		prop.State,
	) {
		return prop.Id, fmt.Errorf("proposal is in wrong state %s: %+v", prop.State.String(), prop)
	}

	return prop.Id, nil
}
