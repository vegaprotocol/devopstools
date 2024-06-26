package market

import (
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/tools"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
)

const MainnetETHUSDTMetadataID = "auto:mainnet_eth_usdt"

var (
	MainnetETHUSDTOracleDecimalPlaces   = uint64(6)
	MainnetETHUSDTMarketSettlementHours = 24 * 31
	MainnetETHUSDTPositionDecimalPlaces = int64(3)
	MainnetETHUSDTDecimalPlaces         = uint64(2)
)

func NewFutureETHUSDTWithoutTime(settlementAsset string, oraclePubKey string) *commandspb.ProposalSubmission {
	return NewFutureETHUSDT(
		settlementAsset,
		oraclePubKey,
		time.Now().Add(30*time.Second),
		time.Now().Add(45*time.Second),
		DefaultExtraMetadata,
	)
}

func NewFutureETHUSDT(
	settlementVegaAssetId string,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	reference := tools.RandAlphaNumericString(40)

	nowTime := time.Now()
	settlementTime := nowTime.Add(time.Hour * (time.Duration(MainnetETHUSDTMarketSettlementHours)))

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title: "ETH/USDT Future",
			Description: `# Summary

This proposal requests to list ETH/USDT-231231 as a market with USDT as a settlement asset on the Vega Network as discussed in: https://community.vega.xyz/.
			
# Rationale

	* ETH is the largest Crypto asset with the highest volume and Marketcap.
	* Given the price, 2 decimal places will be used for price due to the number of valid digits in asset price.
	* Position decimal places will be set to 3 considering the value per contract
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
							Name: "ETH/USDT Future",
							Code: "ETH/USDT-FUTURE",
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
																	Name:                "prices.ETH.value",
																	Type:                datav1.PropertyKey_TYPE_INTEGER,
																	NumberDecimalPlaces: &MainnetETHUSDTOracleDecimalPlaces,
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
										SettlementDataProperty:     "prices.ETH.value",
										TradingTerminationProperty: "vegaprotocol.builtin.timestamp",
									},
								},
							},
						},
						DecimalPlaces: MainnetETHUSDTDecimalPlaces,
						Metadata: append([]string{
							"base:ETH",
							"quote:USDT",
							"class:fx/crypto",
							"quarterly",
							"sector:defi",
							fmt.Sprintf("enactment:%s", enactmentTime.Format(time.RFC3339)),
							fmt.Sprintf("settlement:%s", settlementTime.Format(time.RFC3339)),
							MainnetETHUSDTMetadataID,
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
						PositionDecimalPlaces: MainnetETHUSDTPositionDecimalPlaces,
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
							DisposalSlippageRange: "0.005",
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
