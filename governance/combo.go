package governance

import (
	"context"
	"fmt"

	"github.com/vegaprotocol/devopstools/vegaapi"

	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"code.vegaprotocol.io/vega/wallet/wallet"

	"go.uber.org/zap"
)

func ProposeVoteAndWait(ctx context.Context, description string, proposal *commandspb.ProposalSubmission, proposer wallet.Wallet, proposerPublicKey string, dataNodeClient vegaapi.DataNodeClient, logger *zap.Logger) error {
	proposals := map[string]*commandspb.ProposalSubmission{
		description: proposal,
	}
	return ProposeVoteAndWaitList(ctx, proposals, proposer, proposerPublicKey, dataNodeClient, logger)
}

func ProposeVoteAndWaitList(ctx context.Context, descriptionToProposalConfig map[string]*commandspb.ProposalSubmission, proposer wallet.Wallet, proposerPublicKey string, dataNodeClient vegaapi.DataNodeClient, logger *zap.Logger) error {
	return Propose(ctx, descriptionToProposalConfig, proposer, proposerPublicKey, dataNodeClient, logger)
}

func Propose(ctx context.Context, descriptionToProposalConfig map[string]*commandspb.ProposalSubmission, proposer wallet.Wallet, proposerPublicKey string, dataNodeClient vegaapi.DataNodeClient, logger *zap.Logger) error {
	logger.Info("Submitting proposals", zap.Int("count", len(descriptionToProposalConfig)))
	descriptionToProposalId, err := SubmitProposalList(ctx, descriptionToProposalConfig, proposer, proposerPublicKey, dataNodeClient)
	if err != nil {
		return fmt.Errorf("failed to submit proposals: %w", err)
	}
	logger.Info("Successfully submitted proposals", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))

	logger.Info("Voting for proposals", zap.Int("count", len(descriptionToProposalConfig)))
	if err := VoteOnProposalList(ctx, descriptionToProposalId, proposer, proposerPublicKey, dataNodeClient); err != nil {
		return fmt.Errorf("failed to vote for proposals: %w", err)
	}
	logger.Info("Successfully voted for proposals", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))

	logger.Info("Waiting for proposals enactment", zap.Int("count", len(descriptionToProposalConfig)))
	if err := WaitForEnactList(ctx, descriptionToProposalId, dataNodeClient, logger); err != nil {
		return fmt.Errorf("proposal enactment failed: %w", err)
	}
	logger.Info("Proposals enactment successful", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))

	return nil
}
