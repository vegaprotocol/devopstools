package governance

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
	"golang.org/x/exp/slices"
)

func ProposeVoteProvideLP(
	name string,
	dataNodeClient vegaapi.DataNodeClient,
	lastBlockData *vegaapipb.LastBlockHeightResponse,
	markets []*vega.Market,
	proposerVegawallet *wallet.VegaWallet,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	marketMetadataMarker string,
	proposal *commandspb.ProposalSubmission,
	logger *zap.Logger,
) error {
	market := GetMarket(markets, oraclePubKey, marketMetadataMarker)
	if market != nil {
		logger.Info("market already exist", zap.String("market", market.Id), zap.String("name", name))
		return nil
	}
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
	if err := SubmitTx("propose market", dataNodeClient, proposerVegawallet, logger, &walletTxReq); err != nil {
		return err
	}

	//
	// Find Proposal
	//
	time.Sleep(time.Second * 10)
	res, err := dataNodeClient.ListGovernanceData(&v2.ListGovernanceDataRequest{
		ProposalReference: &reference,
	})
	if err != nil {
		return err
	}
	var proposalId string
	for _, edge := range res.Connection.Edges {
		if edge.Node.Proposal.Reference == reference {
			logger.Info("Found proposal", zap.String("market", name), zap.String("reference", reference),
				zap.String("status", edge.Node.Proposal.State.String()),
				zap.Any("proposal", edge.Node.Proposal))
			proposalId = edge.Node.Proposal.Id
		}
	}

	if len(proposalId) < 1 {
		return fmt.Errorf("got empty proposal id for the %s reference", reference)
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
	if err := SubmitTx("vote on market proposal", dataNodeClient, proposerVegawallet, logger, &voteWalletTxReq); err != nil {
		return err
	}

	//
	// Provide LP
	//

	markets, err = dataNodeClient.GetAllMarkets()
	if err != nil {
		return fmt.Errorf("failed to get markets: %w", err)
	}
	market = GetMarket(markets, oraclePubKey, marketMetadataMarker)
	if market == nil {
		return fmt.Errorf("failed to find particular market for %s oracle public key", oraclePubKey)
	}

	// Prepare vegawallet Transaction Request
	provideLPWalletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_LiquidityProvisionSubmission{
			LiquidityProvisionSubmission: &commandspb.LiquidityProvisionSubmission{
				Fee:              "0.01",
				MarketId:         market.Id,
				CommitmentAmount: "1",
			},
		},
	}
	if err := SubmitTx("Provide LP", dataNodeClient, proposerVegawallet, logger, &provideLPWalletTxReq); err != nil {
		return err
	}

	return nil
}

const OracleAll = "*"

func GetMarket(markets []*vega.Market, oraclePubKey string, metadataTag string) *vega.Market {
	for _, market := range markets {
		if market.TradableInstrument == nil || market.TradableInstrument.Instrument == nil {
			continue
		}

		future := market.TradableInstrument.Instrument.GetFuture()
		perpetual := market.TradableInstrument.Instrument.GetPerpetual()

		if future == nil && perpetual == nil {
			// TODO: Log the unsupported market here?
			continue
		}

		stringSigners := []string{OracleAll}
		if perpetual != nil &&
			perpetual.DataSourceSpecForSettlementData != nil &&
			perpetual.DataSourceSpecForSettlementData.GetData() != nil &&
			perpetual.DataSourceSpecForSettlementData.GetData().GetExternal() != nil {
			external := perpetual.DataSourceSpecForSettlementData.GetData().GetExternal()

			if external.GetOracle() != nil {
				for _, signer := range external.GetOracle().Signers {
					if signer.GetPubKey() != nil {
						stringSigners = append(stringSigners, signer.GetPubKey().GetKey())
					}
					if signer.GetEthAddress() != nil {
						stringSigners = append(stringSigners, signer.GetEthAddress().Address)
					}
				}
			}

			if external.GetEthOracle() != nil {
				stringSigners = append(stringSigners, external.GetEthOracle().Address)
			}
		}

		if future != nil &&
			future.DataSourceSpecForSettlementData != nil &&
			future.DataSourceSpecForSettlementData.GetData() != nil &&
			future.DataSourceSpecForSettlementData.GetData().GetExternal() != nil &&
			future.DataSourceSpecForSettlementData.GetData().GetExternal().GetOracle() != nil {
			signers := future.DataSourceSpecForSettlementData.GetData().GetExternal().GetOracle().Signers
			for _, signer := range signers {
				if signer.GetPubKey() != nil {
					stringSigners = append(stringSigners, signer.GetPubKey().GetKey())
				}
				if signer.GetEthAddress() != nil {
					stringSigners = append(stringSigners, signer.GetEthAddress().Address)
				}
			}
		}

		if slices.Contains(
			stringSigners,
			oraclePubKey,
		) && slices.Contains(
			market.TradableInstrument.Instrument.Metadata.Tags,
			metadataTag,
		) {
			// TODO check if open
			return market
		}
	}
	return nil
}

func TerminateMarketProposal(closingTime, enactmentTime time.Time, marketName string, marketId string, price string) *commandspb.ProposalSubmission {
	reference := tools.RandAlphaNumericString(40)

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       fmt.Sprintf("Terminate %s market", marketName),
			Description: fmt.Sprintf("Terminate %s market", marketName),
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),

			Change: &vega.ProposalTerms_UpdateMarketState{
				UpdateMarketState: &vega.UpdateMarketState{
					Changes: &vega.UpdateMarketStateConfiguration{
						MarketId:   marketId,
						UpdateType: vega.MarketStateUpdateType_MARKET_STATE_UPDATE_TYPE_TERMINATE,
						Price:      &price,
					},
				},
			},
		},
	}
}
