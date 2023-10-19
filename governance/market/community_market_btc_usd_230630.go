package market

import (
	"fmt"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
)

const CommunityBTCUSD230630MetadataID = "auto:community_btcusdc230630"

// settlementVegaAssetId ideally with 6 decimal places
func NewCommunityBTCUSD230630(
	settlementVegaAssetId string,
	decimalPlaces uint64,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) (*commandspb.ProposalSubmission, error) {
	if decimalPlaces != 6 {
		return nil, fmt.Errorf("asset decimal places for market %s(%s) must be 6", "BTC/USDT expiry 2023 June 30th", settlementVegaAssetId)
	}

	nowTime := time.Now()
	return &commandspb.ProposalSubmission{
		Reference: "BTC/USDT-23063",
		Rationale: &vega.ProposalRationale{
			Title:       "BTC/USDT-230630",
			Description: "## Summary\n\nThis proposal requests to list BTC/USDT-230630 as a market with USDT as a settlement asset on the Vega Network as discussed in: https://community.vega.xyz/.\n\n## Rationale\n\n- BTC is the largest Crypto asset with the highest volume and Marketcap.\n- Given the price, 1 decimal places will be used for price due to the number of valid digits in asset price. \n- Position decimal places will be set to 4 considering the value per contract\n- USDT is chosen as settlement asset due to its stability.",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						Instrument: &vega.InstrumentConfiguration{
							Name: "BTC/USDT expiry 2023 June 30th",
							Code: "BTC/USDT-230630",
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
																	Name:                "prices.BTC.value",
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
																	Name: "prices.BTC.timestamp",
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
										SettlementDataProperty:     "prices.BTC.value",
										TradingTerminationProperty: "vegaprotocol.builtin.timestamp",
									},
								},
							},
						},
						DecimalPlaces: 1,
						Metadata: append([]string{
							"base:BTC",
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
						PositionDecimalPlaces: 4,
						LiquiditySlaParameters: &vega.LiquiditySLAParameters{
							PriceRange:                  "0.05",
							CommitmentMinTimeFraction:   "0.95",
							PerformanceHysteresisEpochs: 1,
							SlaCompetitionFactor:        "0.90",
						},
						LinearSlippageFactor:    "0.001",
						QuadraticSlippageFactor: "0.0",
					},
				},
			},
		},
	}, nil
}
