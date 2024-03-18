package governance

import (
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"

	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"

	"go.uber.org/zap"
)

func ProposeVoteAndWait(
	description string,
	proposal *commandspb.ProposalSubmission,
	proposerVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	return ProposeVoteAndWaitList(
		map[string]*commandspb.ProposalSubmission{
			description: proposal,
		},
		proposerVegawallet, dataNodeClient, logger,
	)
}

func ProposeAndVoteList(
	descriptionToProposalConfig map[string]*commandspb.ProposalSubmission,
	proposerVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	return comboList(
		descriptionToProposalConfig,
		proposerVegawallet,
		dataNodeClient,
		logger,
		true,
		false,
	)
}

func ProposeVoteAndWaitList(
	descriptionToProposalConfig map[string]*commandspb.ProposalSubmission,
	proposerVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	return comboList(
		descriptionToProposalConfig,
		proposerVegawallet,
		dataNodeClient,
		logger,
		true,
		true,
	)
}

func comboList(
	descriptionToProposalConfig map[string]*commandspb.ProposalSubmission,
	proposerVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
	vote bool,
	wait bool,
) error {
	//
	// Propose
	//
	logger.Info("Submitting proposals", zap.Int("count", len(descriptionToProposalConfig)))
	descriptionToProposalId, err := SubmitProposalList(
		descriptionToProposalConfig, proposerVegawallet, dataNodeClient, logger,
	)
	if err != nil {
		return fmt.Errorf("failed to submit proposal list, %w", err)
	}
	logger.Info("Successfully submitted proposals", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))
	if vote {
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
	}
	if vote && wait {
		//
		// Wait
		//
		logger.Info("Waiting for Enactment of proposals", zap.Int("count", len(descriptionToProposalConfig)))
		err = WaitForEnactList(
			descriptionToProposalId, dataNodeClient, logger,
		)
		if err != nil {
			return fmt.Errorf("failed to vote on proposal list, %w", err)
		}
		logger.Info("Successfully Enacted all proposals", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))
	}
	return nil
}

func ProposeAndVote(
	logger *zap.Logger,
	proposerVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	proposal *commandspb.ProposalSubmission,
) error {
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
	if err := SubmitTx("submit transaction", dataNodeClient, proposerVegawallet, logger, &walletTxReq); err != nil {
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
	if err := SubmitTx("vote on proposal", dataNodeClient, proposerVegawallet, logger, &voteWalletTxReq); err != nil {
		return fmt.Errorf("failed to vote on proposal %s: %w", proposalId, err)
	}

	return nil
}
