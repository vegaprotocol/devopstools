package governance

import (
	"time"

	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"github.com/vegaprotocol/devopstools/governance/market"
	"github.com/vegaprotocol/devopstools/governance/networkparameters"
)

func MainnetUpgradeBatchProposal(
	closingTime, enactmentTime time.Time,
) *commandspb.BatchProposalSubmission {
	return NewBatchProposal(
		"Update Markets and network parameters to support new v0.74.x features",
		"Update market BTC/USD and ETH/USD and some network parameters",
		time.Now(),
		[]*commandspb.ProposalSubmission{
			networkparameters.NewUpdateParametersProposal(
				"market.liquidity.minimum.probabilityOfTrading.lpOrders", "0.001", closingTime, enactmentTime,
			),
			networkparameters.NewUpdateParametersProposal(
				"market.liquidity.probabilityOfTrading.tau.scaling", "0.1", closingTime, enactmentTime,
			),
			market.NewBTCUSDMainnetMarketProposal(
				market.IncentiveVegaAssetId, closingTime, enactmentTime, []string{},
			),
			market.NewBTCUSDMainnet2MarketProposal(
				market.IncentiveVegaAssetId, closingTime, enactmentTime, []string{},
			),
		},
		nil,
	)
}
