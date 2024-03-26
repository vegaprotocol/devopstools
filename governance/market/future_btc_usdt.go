package market

import (
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/tools"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
)

const MainnetBTCUSDTMetadataID = "auto:mainnet_btc_usdt"

var (
	MainnetBTCUSDTOracleDecimalPlaces   = uint64(6)
	MainnetBTCUSDTMarketSettlementHours = 24 * 31
	MainnetBTCUSDTPositionDecimalPlaces = int64(4)
	MainnetBTCUSDTDecimalPlaces         = uint64(1)
)

func NewFutureBTCUSDTWithoutTime(settlementAsset string, oraclePubKey string) *commandspb.ProposalSubmission {
	return NewFutureBTCUSDT(
		settlementAsset,
		oraclePubKey,
		time.Now().Add(30*time.Second),
		time.Now().Add(45*time.Second),
		DefaultExtraMetadata,
	)
}

// settlementVegaAssetId ideally with 6 decimal places
func NewFutureBTCUSDT(
	settlementVegaAssetId string,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	reference := tools.RandAlphaNumericString(40)

	nowTime := time.Now()
	settlementTime := nowTime.Add(time.Hour * (time.Duration(MainnetBTCUSDTMarketSettlementHours)))

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title: "BTC/USDT Future",
			Description: `# Summary

This proposal requests to list BTC/USDT-231231 as a market with USDT as a settlement asset on the Vega Network as discussed in: https://community.vega.xyz/.
			
# Rationale
			
	* BTC is the largest Crypto asset with the highest volume and Marketcap.
	* Given the price, 1 decimal places will be used for price due to the number of valid digits in asset price.
	* Position decimal places will be set to 4 considering the value per contract
	* USDT is chosen as settlement asset due to its stability.`,
		},

		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						TickSize: "1",
						Instrument: &vega.InstrumentConfiguration{
							Name: "BTC/USDT FUTURE",
							Code: "BTC/USDT-MAINNET",
							Product: &vega.InstrumentConfiguration_Future{
								Future: &vega.FutureProduct{
									SettlementAsset: settlementVegaAssetId,
									QuoteName:       "USDT",
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
																	NumberDecimalPlaces: &MainnetBTCUSDTOracleDecimalPlaces,
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
																		Value:    fmt.Sprintf("%d", settlementTime.Unix()),
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
																Value:    fmt.Sprintf("%d", settlementTime.Unix()),
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
						DecimalPlaces: MainnetBTCUSDTDecimalPlaces,
						Metadata: append([]string{
							"base:BTC",
							"quote:USDT",
							"class:fx/crypto",
							"quarterly",
							"sector:defi",
							fmt.Sprintf("enactment:%s", enactmentTime.Format(time.RFC3339)),
							fmt.Sprintf("settlement:%s", settlementTime.Format(time.RFC3339)),
							MainnetBTCUSDTMetadataID,
						}, extraMetadata...),
						PriceMonitoringParameters: &vega.PriceMonitoringParameters{
							Triggers: []*vega.PriceMonitoringTrigger{
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
								Tau:                   0.000009506426342,
								Params: &vega.LogNormalModelParams{
									Mu:    0,
									R:     0,
									Sigma: 1.5,
								},
							},
						},
						PositionDecimalPlaces: MainnetBTCUSDTPositionDecimalPlaces,
						LpPriceRange:          &[]string{"0.05"}[0],
						LiquiditySlaParameters: &vega.LiquiditySLAParameters{
							PerformanceHysteresisEpochs: 1,
							PriceRange:                  "0.05",
							SlaCompetitionFactor:        "0.90",
							CommitmentMinTimeFraction:   "0.95",
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
						QuadraticSlippageFactor: "0",
					},
				},
			},
		},
	}
}
