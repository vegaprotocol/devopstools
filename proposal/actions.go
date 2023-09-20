package proposal

import (
	"fmt"
	"time"

	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
)

func ProposeAndVoteList(
	descriptionToProposalConfig map[string]*commandspb.ProposalSubmission,
	proposerVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	logger.Info("Submitting proposals", zap.Int("count", len(descriptionToProposalConfig)))
	descriptionToProposalId, err := SubmitProposalList(
		descriptionToProposalConfig, proposerVegawallet, dataNodeClient, logger,
	)
	if err != nil {
		return fmt.Errorf("failed to submit proposal list, %w", err)
	}
	logger.Info("Successfully submitted proposals", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))
	//
	// Vote
	//
	logger.Info("Voting on proposals", zap.Int("count", len(descriptionToProposalConfig)))
	err = VoteOnProposalList(
		descriptionToProposalId, proposerVegawallet, dataNodeClient, logger,
	)
	if err != nil {
		return fmt.Errorf("failed to vote on proposal list, %w", err)
	}
	logger.Info("Successfully voted on proposals", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))
	return nil
}

func ProposeAndVote(
	logger *zap.Logger,
	proposerVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	proposal *commandspb.ProposalSubmission) error {

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
	if err := submitTx("submit transaction", dataNodeClient, proposerVegawallet, logger, &walletTxReq); err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	//
	// Find proposal
	//
	proposalId, err := tools.RetryReturn[string](10, 6*time.Second, func() (string, error) {
		logger.Info("Looking for proposal", zap.String("reference", reference))

		res, err := dataNodeClient.ListGovernanceData(&v2.ListGovernanceDataRequest{
			ProposalReference: &reference,
		})
		if err != nil {
			return "", fmt.Errorf("failed to list governance proposals: %w", err)
		}
		var proposalId string
		for _, edge := range res.Connection.Edges {
			if edge.Node.Proposal.Reference == reference {
				logger.Info("Found proposal", zap.String("reference", reference),
					zap.String("status", edge.Node.Proposal.State.String()),
					zap.Any("proposal", edge.Node.Proposal))

				if edge.Node.Proposal.ErrorDetails != nil && len(*edge.Node.Proposal.ErrorDetails) > 0 {
					return "", fmt.Errorf("proposal failed: %s", *edge.Node.Proposal.ErrorDetails)
				}
				proposalId = edge.Node.Proposal.Id
			}
		}

		if len(proposalId) < 1 {
			return "", fmt.Errorf("proposal not found")
		}

		return proposalId, nil
	})

	if err != nil {
		return fmt.Errorf("failed to get proposal ID for reference: %s: %w", reference, err)
	}

	//
	// VOTE
	//
	voteWalletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
			VoteSubmission: &commandspb.VoteSubmission{
				ProposalId: proposalId,
				Value:      vega.Vote_VALUE_YES,
			},
		},
	}
	if err := submitTx("vote on proposal", dataNodeClient, proposerVegawallet, logger, &voteWalletTxReq); err != nil {
		return fmt.Errorf("failed to vote on proposal %s: %w", proposalId, err)
	}

	return nil
}

func submitTx(
	description string,
	dataNodeClient vegaapi.DataNodeClient,
	proposerVegawallet *wallet.VegaWallet,
	logger *zap.Logger,
	walletTxReq *walletpb.SubmitTransactionRequest,
) error {
	lastBlockData, err := dataNodeClient.LastBlockData()
	if err != nil {
		return fmt.Errorf("failed to submit tx: %w", err)
	}

	// Sign + Proof of Work vegawallet Transaction request
	signedTx, err := proposerVegawallet.SignTxWithPoW(walletTxReq, lastBlockData)
	if err != nil {
		logger.Error("Failed to sign a trasnaction", zap.String("description", description),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", &walletTxReq), zap.Error(err))
		return err
	}

	// wrap in vega Transaction Request
	submitReq := &vegaapipb.SubmitTransactionRequest{
		Tx:   signedTx,
		Type: vegaapipb.SubmitTransactionRequest_TYPE_SYNC,
	}

	// Submit Transaction
	logger.Info("Submit transaction", zap.String("description", description),
		zap.String("proposer", proposerVegawallet.PublicKey))
	submitResponse, err := dataNodeClient.SubmitTransaction(submitReq)
	if err != nil {
		logger.Error("Failed to submit a trasnaction", zap.String("description", description),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", submitReq), zap.Error(err))
		return err
	}
	if !submitResponse.Success {
		logger.Error("Transaction submission response is not successful",
			zap.String("proposer", proposerVegawallet.PublicKey), zap.String("description", description),
			zap.Any("txReq", submitReq.String()), zap.String("response", fmt.Sprintf("%#v", submitResponse)))
		return err
	}
	logger.Info("Successful Submision of Market Proposal", zap.String("description", description),
		zap.String("proposer", proposerVegawallet.PublicKey), zap.String("txHash", submitResponse.TxHash),
		zap.Any("response", submitResponse))

	return nil
}
