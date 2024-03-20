package market

import (
	"fmt"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
)

const CommunityETHUSD230630MetadataID = "auto:community_ethusdc230630"

// settlementVegaAssetId ideally with 6 decimal places
func NewCommunityETHUSD230630(
	settlementVegaAssetId string,
	decimalPlaces uint64,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) (*commandspb.ProposalSubmission, error) {
	if decimalPlaces != 6 {
		return nil, fmt.Errorf("asset decimal places for market %s(%s) must be 6", "ETH/USDT expiry 2023 June 30th", settlementVegaAssetId)
	}

	nowTime := time.Now()
	return &commandspb.ProposalSubmission{
		Reference: "ETH/USDT-23063",
		Rationale: &vega.ProposalRationale{
			Title:       "ETH/USDT-230630",
			Description: "## Summary\n\nThis proposal requests to list ETH/USDT-230630 as a market with USDT as a settlement asset on the Vega Network as discussed in: https://community.vega.xyz/.\n\n## Rationale\n\n- ETH is the second largest Crypto asset with the highest volume and Marketcap.\n- Given the price, 2 decimal places will be used for price due to the number of valid digits in asset price. \n- Position decimal places will be set to 3 considering the value per contract\n- USDT is chosen as settlement asset due to its stability.",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						TickSize: "1",
						Instrument: &vega.InstrumentConfiguration{
							Name: "ETH/USDT expiry 2023 June 30th",
							Code: "ETH/USDT-230630",
							Product: &vega.InstrumentConfiguration_Future{
								Future: &vega.FutureProduct{
									SettlementAsset: settlementVegaAssetId,
									QuoteName:       "USDC",
									DataSourceSpecForSettlementData: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_External{
											External: &vega.DataSourceDefinitionExternal{
												SourceType: &vega.DataSourceDefinitionExternal_Oracle{
													Oracle: &vega.DataSourceSpecConfiguration{
														Signers: []*datav1.Signer{
															{
																Signer: &datav1.Signer_EthAddress{
																	EthAddress: &datav1.ETHAddress{
																		Address: oraclePubKey,
																	},
																},
															},
														},
														Filters: []*datav1.Filter{
															{
																Key: &datav1.PropertyKey{
																	Name:                "prices.ETH.value",
																	Type:                datav1.PropertyKey_TYPE_INTEGER,
																	NumberDecimalPlaces: &decimalPlaces,
																},
																Conditions: []*datav1.Condition{
																	{
																		Operator: datav1.Condition_OPERATOR_GREATER_THAN,
																		Value:    "0",
																	},
																},
															},
															{
																Key: &datav1.PropertyKey{
																	Name: "prices.ETH.timestamp",
																	Type: datav1.PropertyKey_TYPE_TIMESTAMP,
																},
																Conditions: []*datav1.Condition{
																	{
																		Operator: datav1.Condition_OPERATOR_GREATER_THAN_OR_EQUAL,
																		Value:    fmt.Sprintf("%d", nowTime.Add(time.Hour*(24*31+1)).Unix()),
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
										SourceType: &vega.DataSourceDefinition_Internal{
											Internal: &vega.DataSourceDefinitionInternal{
												SourceType: &vega.DataSourceDefinitionInternal_Time{
													Time: &vega.DataSourceSpecConfigurationTime{
														Conditions: []*datav1.Condition{
															{
																Operator: datav1.Condition_OPERATOR_GREATER_THAN_OR_EQUAL,
																Value:    fmt.Sprintf("%d", nowTime.Add(time.Hour*24*31).Unix()),
															},
														},
													},
												},
											},
										},
									},
									DataSourceSpecBinding: &vega.DataSourceSpecToFutureBinding{
										SettlementDataProperty:     "prices.ETH.value",
										TradingTerminationProperty: "vegaprotocol.builtin.timestamp",
									},
								},
							},
						},
						DecimalPlaces: 2,
						Metadata: append([]string{
							"base:ETH",
							"quote:USD",
							"class:fx/crypto",
							"quarterly",
							"sector:defi",
							"managed:vega/ops",
							"enactment:2023-05-21T08:00:00Z",
							"settlement:2023-06-30T08:00:00Z",
						}, extraMetadata...),
						PriceMonitoringParameters: &vega.PriceMonitoringParameters{
							Triggers: []*vega.PriceMonitoringTrigger{
								{
									Horizon:          60,
									Probability:      "0.9999",
									AuctionExtension: 5,
								},
								{
									Horizon:          600,
									Probability:      "0.9999",
									AuctionExtension: 30,
								},
								{
									Horizon:          3600,
									Probability:      "0.9999",
									AuctionExtension: 120,
								},
								{
									Horizon:          14400,
									Probability:      "0.9999",
									AuctionExtension: 180,
								},
								{
									Horizon:          43200,
									Probability:      "0.9999",
									AuctionExtension: 300,
								},
							},
						},
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 1,
							},
							TriggeringRatio:  "0.7",
							AuctionExtension: 1,
						},

						RiskParameters: &vega.NewMarketConfiguration_LogNormal{
							LogNormal: &vega.LogNormalRiskModel{
								RiskAversionParameter: 0.000001,
								Tau:                   0.0001140771161,
								Params: &vega.LogNormalModelParams{
									Sigma: 1.5,
								},
							},
						},
						PositionDecimalPlaces: 3,
						LiquiditySlaParameters: &vega.LiquiditySLAParameters{
							PriceRange:                  "0.05",
							CommitmentMinTimeFraction:   "0.95",
							PerformanceHysteresisEpochs: 1,
							SlaCompetitionFactor:        "0.90",
						},
						LiquidationStrategy: &vega.LiquidationStrategy{
							DisposalTimeStep:    30,
							MaxFractionConsumed: "0.1",
							DisposalFraction:    "0.1",
							// FullDisposalSize:    0,
						},
						MarkPriceConfiguration: &vega.CompositePriceConfiguration{
							CompositePriceType: vega.CompositePriceType_COMPOSITE_PRICE_TYPE_LAST_TRADE,
						},
						LinearSlippageFactor:    "0.001",
						QuadraticSlippageFactor: "0.0",
					},
				},
			},
		},
	}, nil
}
