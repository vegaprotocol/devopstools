package market

import (
	"fmt"
	"time"

	dstypes "code.vegaprotocol.io/vega/core/datasource/common"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

func NewETHBTCMarketProposal(
	settlementVegaAssetId string,
	decimalPlaces uint64,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	var (
		reference = tools.RandAlphaNumericString(40)
		Name      = fmt.Sprintf("ETHBTC Quarterly (%s)", time.Now().AddDate(0, 3, 0).Format("Jan 2006")) // Now + 3 months
		pubKey    = dstypes.CreateSignerFromString(oraclePubKey, dstypes.SignerTypePubKey)
	)

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "New BTC market",
			Description: "New BTC market",
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
						QuadraticSlippageFactor: "0.0",
						Instrument: &vega.InstrumentConfiguration{
							Name: Name,
							Code: "ETHBTC.QM21",
							Product: &vega.InstrumentConfiguration_Future{
								Future: &vega.FutureProduct{
									SettlementAsset: settlementVegaAssetId,
									QuoteName:       "BTC",
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
							"formerly:1F0BB6EB5703B099",
							"base:ETH",
							"quote:BTC",
							"class:fx/crypto",
							"quarterly",
							"managed:vega/ops",
							"sector:crypto",
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
						LiquiditySlaParameters: &vega.LiquiditySLAParameters{
							PriceRange:                  "0.05",
							CommitmentMinTimeFraction:   "0.95",
							PerformanceHysteresisEpochs: 1,
							SlaCompetitionFactor:        "0.90",
						},
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 10,
							},
							TriggeringRatio:  "0.7",
							AuctionExtension: 1,
						},
						LiquidationStrategy: &vega.LiquidationStrategy{
							DisposalTimeStep:    30,
							MaxFractionConsumed: "0.1",
							DisposalFraction:    "0.1",
							// FullDisposalSize:    0,
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
						MarkPriceConfiguration: &vega.CompositePriceConfiguration{
							CompositePriceType: vega.CompositePriceType_COMPOSITE_PRICE_TYPE_LAST_TRADE,
						},
					},
				},
			},
		},
	}
}
