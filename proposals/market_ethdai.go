package proposals

import (
	"fmt"
	"time"

	"code.vegaprotocol.io/vega/core/types"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

func NewETHDAIMarketProposal(
	settlementVegaAssetId string,
	decimalPlaces uint64,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	var (
		reference = tools.RandAlpaNumericString(40)
		Name      = fmt.Sprintf("ETHDAI Monthly (%s)", time.Now().AddDate(0, 1, 0).Format("Jan 2006")) // Now + 1 month
		pubKey    = types.CreateSignerFromString(oraclePubKey, types.DataSignerTypePubKey)
	)

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "New ETHDAI market",
			Description: "New ETHDAI market",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						DecimalPlaces:           decimalPlaces,
						PositionDecimalPlaces:   3,
						LinearSlippageFactor:    "0.1",
						QuadraticSlippageFactor: "0.1",
						Instrument: &vega.InstrumentConfiguration{
							Name: Name,
							Code: "ETHDAI.MF21",
							Product: &vega.InstrumentConfiguration_Future{
								Future: &vega.FutureProduct{
									SettlementAsset: settlementVegaAssetId,
									QuoteName:       "DAI",
									DataSourceSpecForSettlementData: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_External{
											External: &vega.DataSourceDefinitionExternal{
												SourceType: &vega.DataSourceDefinitionExternal_Oracle{
													Oracle: &vega.DataSourceSpecConfiguration{
														Signers: []*datav1.Signer{pubKey.IntoProto()},
														Filters: []*datav1.Filter{
															{
																Key: &datav1.PropertyKey{
																	Name: "prices.ETH.value",
																	Type: datav1.PropertyKey_TYPE_INTEGER,
																},
																Conditions: []*datav1.Condition{
																	{
																		Operator: datav1.Condition_OPERATOR_EQUALS,
																		Value:    "1",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									DataSourceSpecForTradingTermination: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_External{
											External: &vega.DataSourceDefinitionExternal{
												SourceType: &vega.DataSourceDefinitionExternal_Oracle{
													Oracle: &vega.DataSourceSpecConfiguration{
														Signers: []*datav1.Signer{pubKey.IntoProto()},
														Filters: []*datav1.Filter{
															{
																Key: &datav1.PropertyKey{
																	Name: "termination.ETH.value",
																	Type: datav1.PropertyKey_TYPE_BOOLEAN,
																},
																Conditions: []*datav1.Condition{
																	{
																		Operator: datav1.Condition_OPERATOR_EQUALS,
																		Value:    "1",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									DataSourceSpecBinding: &vega.DataSourceSpecToFutureBinding{
										SettlementDataProperty:     "prices.ETH.value",
										TradingTerminationProperty: "termination.ETH.value",
									},
								},
							},
						},
						Metadata: append([]string{
							"formerly:4BC6D2154BE74E1F",
							"base:ETH",
							"quote:DAI",
							"class:fx/crypto",
							"monthly",
							"sector:defi",
						}, extraMetadata...),
						PriceMonitoringParameters: &vega.PriceMonitoringParameters{
							Triggers: []*vega.PriceMonitoringTrigger{
								{
									Horizon:          43200,
									Probability:      "0.9999999",
									AuctionExtension: 600,
								},
							},
						},
						LpPriceRange: "0.5",
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 10,
							},
							TriggeringRatio:  "0.7",
							AuctionExtension: 1,
						},
						RiskParameters: &vega.NewMarketConfiguration_LogNormal{
							LogNormal: &vega.LogNormalRiskModel{
								RiskAversionParameter: 0.01,
								Tau:                   0.0001140771161,
								Params: &vega.LogNormalModelParams{
									Mu:    0,
									R:     0.016,
									Sigma: 0.3,
								},
							},
						},
					},
				},
			},
		},
	}
}
