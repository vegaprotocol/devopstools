package governance

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"

	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"

	"go.uber.org/zap"
)

func ProposeVoteAndWait(
	description string,
	proposal *commandspb.ProposalSubmission,
	proposer *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	return ProposeVoteAndWaitList(
		map[string]*commandspb.ProposalSubmission{
			description: proposal,
		},
		proposer, dataNodeClient, logger,
	)
}

func ProposeVoteAndWaitList(
	descriptionToProposalConfig map[string]*commandspb.ProposalSubmission,
	proposer *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	return Propose(
		descriptionToProposalConfig,
		proposer,
		dataNodeClient,
		logger,
		true,
		true,
	)
}

func Propose(
	descriptionToProposalConfig map[string]*commandspb.ProposalSubmission,
	proposer *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
	vote bool,
	wait bool,
) error {
	logger.Info("Submitting proposals", zap.Int("count", len(descriptionToProposalConfig)))
	descriptionToProposalId, err := SubmitProposalList(descriptionToProposalConfig, proposer, dataNodeClient, logger)
	if err != nil {
		return fmt.Errorf("failed to submit proposals: %w", err)
	}
	logger.Info("Successfully submitted proposals", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))

	if vote {
		logger.Info("Voting for proposals", zap.Int("count", len(descriptionToProposalConfig)))
		if err := VoteOnProposalList(descriptionToProposalId, proposer, dataNodeClient, logger); err != nil {
			return fmt.Errorf("failed to vote for proposals: %w", err)
		}
		logger.Info("Successfully voted for proposals", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))
	}

	if vote && wait {
		logger.Info("Waiting for proposals enactment", zap.Int("count", len(descriptionToProposalConfig)))
		if err := WaitForEnactList(descriptionToProposalId, dataNodeClient, logger); err != nil {
			return fmt.Errorf("proposal enactment failed: %w", err)
		}
		logger.Info("Proposals enactment successful", zap.Int("count", len(descriptionToProposalId)), zap.Any("details", descriptionToProposalId))
	}

	return nil
}
