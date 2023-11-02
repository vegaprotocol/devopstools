package propose

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"code.vegaprotocol.io/vega/core/netparams"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/governance"
	assetsgov "github.com/vegaprotocol/devopstools/governance/assets"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"go.uber.org/zap"
)

type ProposeAssetArgs struct {
	*ProposeArgs
}

var proposeAssetArgs ProposeAssetArgs

// provideLPCmd represents the provideLP command
var proposeAssetCmd = &cobra.Command{
	Use:   "assets",
	Short: "Propose asset",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ProposeAssetRun(&proposeAssetArgs); err != nil {
			proposeAssetArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	proposeAssetArgs.ProposeArgs = &proposeArgs
	ProposeCmd.AddCommand(proposeAssetCmd)
}

func getAssets(env string) []assetsgov.AssetProposalDetails {
	result := []assetsgov.AssetProposalDetails{}

	if env != types.NetworkFairground {
		return result
	}

	// Fairground Assets
	result = []assetsgov.AssetProposalDetails{
		{
			Name:                     "Wrapped Ether",
			Symbol:                   "WETH",
			Decimals:                 18,
			Quantum:                  big.NewInt(500000000000000),
			ERC20Address:             "0xC0DF5eB3e58f21026E6F997Dd9C3e1Fa07d22428",
			ERC20LifetimeLimit:       "50000000000000000000",
			ERC20WithdrawalThreshold: "1",
		},
		{
			Name:                     "Dai Stablecoin",
			Symbol:                   "DAI",
			Decimals:                 18,
			Quantum:                  big.NewInt(1000000000000000000),
			ERC20Address:             "0x12ba0B32016099d19E685113123ba7018d702EA2",
			ERC20LifetimeLimit:       "10000000000000000000000",
			ERC20WithdrawalThreshold: "1",
		},
		{
			Name:                     "Wrapped Ether",
			Symbol:                   "WETH",
			Decimals:                 18,
			Quantum:                  big.NewInt(500000000000000),
			ERC20Address:             "0xC0DF5eB3e58f21026E6F997Dd9C3e1Fa07d22428",
			ERC20LifetimeLimit:       "50000000000000000000",
			ERC20WithdrawalThreshold: "1",
		},
		{
			Name:                     "Tether USD",
			Symbol:                   "USDT",
			Decimals:                 6,
			Quantum:                  big.NewInt(1000000),
			ERC20Address:             "0xdb10bf403771e44d0456f6c51ee655bb67ab05d9",
			ERC20LifetimeLimit:       "100000000000",
			ERC20WithdrawalThreshold: "1",
		},
		{
			Name:                     "USDC",
			Symbol:                   "USDC",
			Decimals:                 6,
			Quantum:                  big.NewInt(1000000),
			ERC20Address:             "0x3af2115f79614e66d98560ded281693909499192",
			ERC20LifetimeLimit:       "100000000000",
			ERC20WithdrawalThreshold: "1",
		},
	}

	return result
}

func isAssetSame(deployed *vega.AssetDetails, proposed *assetsgov.AssetProposalDetails) bool {
	return false
}

func ProposeAssetRun(args *ProposeAssetArgs) error {
	logger := args.Logger
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

	networkAssets, err := network.DataNodeClient.GetAssets()
	if err != nil {
		return fmt.Errorf("failed to get assets already created on the network: %w", err)
	}

	assetsByERC20 := map[string]*vega.AssetDetails{}
	for _, asset := range networkAssets {
		erc20 := asset.GetErc20()
		if erc20 == nil {
			continue
		}
		assetsByERC20[erc20.ContractAddress] = asset
	}

	var errs *multierror.Error

	for _, assetDetails := range getAssets(network.Network) {
		logger.Info("Proposing ERC20 asset",
			zap.String("name", assetDetails.Name),
			zap.String("symbol", assetDetails.Symbol),
			zap.Uint64("decimals", assetDetails.Decimals),
			zap.String("quantum", assetDetails.Quantum.String()),
			zap.String("erc20-address", assetDetails.ERC20Address),
			zap.String("withdrawal-threshold", assetDetails.ERC20WithdrawalThreshold),
			zap.String("lifetime-limit", assetDetails.ERC20LifetimeLimit),
		)

		if _, assetAlreadyCreated := assetsByERC20[assetDetails.ERC20Address]; assetAlreadyCreated {
			logger.Info("Asset with given erc20 address already created",
				zap.String("name", assetDetails.Name),
				zap.String("symbol", assetDetails.Symbol),
				zap.String("erc20-address", assetDetails.ERC20Address),
			)

			continue
		}

		minClose, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinClose])
		if err != nil {
			cError := fmt.Errorf("failed to parse duration for %s: %w", netparams.GovernanceProposalMarketMinClose, err)
			logger.Error(
				"failed to propose asset",
				zap.Error(cError),
				zap.String("name", assetDetails.Name),
			)
			errs = multierror.Append(cError)
			continue
		}
		minEnact, err := time.ParseDuration(network.NetworkParams.Params[netparams.GovernanceProposalMarketMinEnact])
		if err != nil {
			cError := fmt.Errorf("failed to parse duration for %s: %w", netparams.GovernanceProposalMarketMinEnact, err)
			logger.Error(
				"failed to propose asset",
				zap.Error(cError),
				zap.String("name", assetDetails.Name),
			)
			errs = multierror.Append(cError)
			continue
		}

		closingTime := time.Now().Add(time.Second * 20).Add(minClose)
		enactmentTime := time.Now().Add(time.Second * 30).Add(minClose).Add(minEnact)

		proposal := assetsgov.NewAssetProposal(closingTime, enactmentTime, assetDetails)

		args.Logger.Info("Proposing asset: Sending proposal", zap.String("name", assetDetails.Name))
		walletTxReq := walletpb.SubmitTransactionRequest{
			PubKey: network.VegaTokenWhale.PublicKey,
			Command: &walletpb.SubmitTransactionRequest_ProposalSubmission{
				ProposalSubmission: proposal,
			},
		}

		if err := governance.SubmitTx("propose asset", network.DataNodeClient, network.VegaTokenWhale, args.Logger, &walletTxReq); err != nil {
			cError := fmt.Errorf("failed to submit propose asset transaction: %w", err)
			logger.Error(
				"failed to propose asset",
				zap.Error(cError),
				zap.String("name", assetDetails.Name),
			)
			errs = multierror.Append(cError)
			continue
		}

		args.Logger.Info(
			"Proposing asset: Waiting for proposal to be picked up on the network",
			zap.String("name", assetDetails.Name),
			zap.String("ref", proposal.Reference),
		)
		proposalId, err := tools.RetryReturn(6, 10*time.Second, func() (string, error) {
			reference := proposal.Reference

			res, err := network.DataNodeClient.ListGovernanceData(&v2.ListGovernanceDataRequest{
				ProposalReference: &reference,
			})
			if err != nil {
				return "", fmt.Errorf("failed to list governance data for reference %s: %w", proposal.Reference, err)
			}
			var proposalId string
			for _, edge := range res.Connection.Edges {
				if edge.Node.Proposal.Reference == reference {
					args.Logger.Info("Found proposal", zap.String("reference", reference),
						zap.String("status", edge.Node.Proposal.State.String()),
						zap.Any("proposal", edge.Node.Proposal),
						zap.String("name", assetDetails.Name),
					)
					proposalId = edge.Node.Proposal.Id
					break
				}
			}

			if len(proposalId) < 1 {
				return "", fmt.Errorf("got empty proposal id for the %s reference", reference)
			}

			return proposalId, nil
		})

		if err != nil {
			cError := fmt.Errorf("failed to find proposal for new asset:%w", err)

			logger.Error(
				"failed to propose asset",
				zap.Error(cError),
				zap.String("name", assetDetails.Name),
			)
			errs = multierror.Append(cError)
			continue
		}

		args.Logger.Info("Proposing asset: Proposal found",
			zap.String("proposal-id", proposalId),
			zap.String("name", assetDetails.Name),
		)

		args.Logger.Info("Proposing asset: Voting on proposal",
			zap.String("proposal-id", proposalId),
			zap.String("name", assetDetails.Name),
		)

		voteWalletTxReq := walletpb.SubmitTransactionRequest{
			PubKey: network.VegaTokenWhale.PublicKey,
			Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
				VoteSubmission: &commandspb.VoteSubmission{
					ProposalId: proposalId,
					Value:      vega.Vote_VALUE_YES,
				},
			},
		}
		if err := governance.SubmitTx("vote on asset proposal", network.DataNodeClient, network.VegaTokenWhale, args.Logger, &voteWalletTxReq); err != nil {
			cError := fmt.Errorf("failet to vote on the asset proposal: %w")

			logger.Error(
				"failed to propose asset",
				zap.Error(cError),
				zap.String("name", assetDetails.Name),
			)
			errs = multierror.Append(cError)
			continue
		}

		args.Logger.Info("Proposing asset: Voted on proposal",
			zap.String("proposal-id", proposalId),
			zap.String("name", assetDetails.Name),
		)
	}

	return errs.ErrorOrNil()
}
