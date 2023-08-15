package proposals

import (
	"time"

	"code.vegaprotocol.io/vega/libs/ptr"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

const PerpetualBTCUSD = "auto:perpetual_btc_usd"
const PerpetualBTCUSDOracleAddress = "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43"

func NewBTCUSDPerpetualMarketProposal(
	settlementVegaAssetId string,
	decimalPlaces uint64,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	var (
		reference = tools.RandAlpaNumericString(40)
		name      = "BTCUSD Perpetual"
	)

	contractABI := `[{"inputs":[],"name":"latestAnswer","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"}]`

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "New BTCUSD perpetual market",
			Description: "New BTCUSD perpetual Market",
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
							Name: name,
							Code: "BTCUSD.MF21",
							Product: &vega.InstrumentConfiguration_Perpetual{
								Perpetual: &vega.PerpetualProduct{
									ClampLowerBound:     "0",
									ClampUpperBound:     "0",
									InterestRate:        "0",
									MarginFundingFactor: "0.1",
									SettlementAsset:     settlementVegaAssetId,
									QuoteName:           "USD",
									DataSourceSpecForSettlementData: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_External{
											External: &vega.DataSourceDefinitionExternal{
												SourceType: &vega.DataSourceDefinitionExternal_EthOracle{
													EthOracle: &vega.EthCallSpec{
														// https://docs.chain.link/data-feeds/price-feeds/addresses#Sepolia%20Testnet
														Address: oraclePubKey, // chainlink BTC/USD
														Abi:     contractABI,
														Method:  "latestAnswer",
														Normalisers: []*vega.Normaliser{
															{
																Name:       "btc.price",
																Expression: "$[0]",
															},
														},
														RequiredConfirmations: 1,
														Trigger: &vega.EthCallTrigger{
															Trigger: &vega.EthCallTrigger_TimeTrigger{
																TimeTrigger: &vega.EthTimeTrigger{
																	Every: ptr.From(uint64(5)),
																},
															},
														},
														Filters: []*datav1.Filter{
															{
																Key: &datav1.PropertyKey{
																	Name:                "btc.price",
																	Type:                datav1.PropertyKey_TYPE_INTEGER,
																	NumberDecimalPlaces: ptr.From(uint64(8)),
																},
																Conditions: []*datav1.Condition{
																	{
																		Operator: datav1.Condition_OPERATOR_GREATER_THAN_OR_EQUAL,
																		Value:    "0",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									DataSourceSpecForSettlementSchedule: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_Internal{
											Internal: &vega.DataSourceDefinitionInternal{
												SourceType: &vega.DataSourceDefinitionInternal_TimeTrigger{
													TimeTrigger: &vega.DataSourceSpecConfigurationTimeTrigger{
														Conditions: []*datav1.Condition{
															{
																Operator: datav1.Condition_OPERATOR_GREATER_THAN_OR_EQUAL,
																Value:    "0",
															},
														},
														Triggers: []*datav1.InternalTimeTrigger{
															{
																Every: 300, // 5 mins in seconds
															},
														},
													},
												},
											},
										},
									},
									DataSourceSpecBinding: &vega.DataSourceSpecToPerpetualBinding{
										SettlementDataProperty:     "btc.price",
										SettlementScheduleProperty: "vegaprotocol.builtin.timetrigger",
									},
								},
							},
						},
						Metadata: append([]string{
							"formerly:50657270657475616c",
							"base:BTC",
							"quote:USD",
							"class:fx/crypto",
							"monthly",
							"sector:crypto",
						}, extraMetadata...),
						PriceMonitoringParameters: &vega.PriceMonitoringParameters{
							Triggers: []*vega.PriceMonitoringTrigger{
								{
									Horizon:          43200,
									Probability:      "0.9999999",
									AuctionExtension: 600,
								},
								{
									Horizon:          300,
									Probability:      "0.9999",
									AuctionExtension: 60,
								},
							},
						},
						LpPriceRange: "0.5",
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 10,
							},
							TriggeringRatio:  "0.0",
							AuctionExtension: 1,
						},
						RiskParameters: &vega.NewMarketConfiguration_LogNormal{
							LogNormal: &vega.LogNormalRiskModel{
								RiskAversionParameter: 0.0001,
								Tau:                   0.0000190129,
								Params: &vega.LogNormalModelParams{
									Mu:    0,
									R:     0.016,
									Sigma: 1.25,
								},
							},
						},
					},
				},
			},
		},
	}
}

// "proposalSubmission": {
// 	"reference": "injected_at_runtime",
// 	"rationale": {
// 	  "title": "BTCUSD market",
// 	  "description": "New BTCUSD market"
// 	},
// 	"terms": {
// 	  "closingTimestamp": 0,
// 	  "enactmentTimestamp": 0,
// 	  "newMarket": {
// 		"changes": {
// 		  "instrument": {
// 			"name": "BTCUSD Monthly (30 Jun 2022)",
// 			"code": "BTCUSD.MF21",
// 			"future": {
// 			  "settlementAsset": "fDAI",
// 			  "quoteName": "USD",
// 			  "oracleSpecForSettlementPrice": {
// 				"pubKeys": [
// 				  "0x51a3a77554709b5db7b769d4376560ad6398c7b08380b5a3e49bda1236697f4f"
// 				],
// 				"filters": [
// 				  {
// 					"key": {
// 					  "name": "prices.BTC.value",
// 					  "type": "TYPE_INTEGER"
// 					},
// 					"conditions": [
// 					  {
// 						"operator": "OPERATOR_EQUALS",
// 						"value": "1"
// 					  }
// 					]
// 				  }
// 				]
// 			  },
// 			  "oracleSpecForTradingTermination": {
// 				"pubKeys": [
// 				  "0x51a3a77554709b5db7b769d4376560ad6398c7b08380b5a3e49bda1236697f4f"
// 				],
// 				"filters": [
// 				  {
// 					"key": {
// 					  "name": "termination.BTC.value",
// 					  "type": "TYPE_BOOLEAN"
// 					},
// 					"conditions": [
// 					  {
// 						"operator": "OPERATOR_EQUALS",
// 						"value": "1"
// 					  }
// 					]
// 				  }
// 				]
// 			  },
// 			  "oracleSpecBinding": {
// 				"settlementPriceProperty": "prices.BTC.value",
// 				"tradingTerminationProperty": "termination.BTC.value"
// 			  }
// 			}
// 		  },
// 		  "decimalPlaces": 5,
// 		  "positionDecimalPlaces": 3,
// 		  "metadata": [
// 			"formerly:076BB86A5AA41E3E",
// 			"base:BTC",
// 			"quote:USD",
// 			"class:fx/crypto",
// 			"monthly",
// 			"sector:crypto"
// 		  ],
// 		  "priceMonitoringParameters": {
// 			"triggers": [
// 			  {
// 				"horizon": 43200,
// 				"probability": "0.9999999",
// 				"auctionExtension": 600
// 			  },
// 			  {
// 				"horizon": 300,
// 				"probability": "0.9999",
// 				"auctionExtension": 60
// 			  }
// 			]
// 		  },
// 		  "liquidityMonitoringParameters": {
// 			"targetStakeParameters": {
// 			  "timeWindow": 3600,
// 			  "scalingFactor": 10
// 			},
// 			"triggeringRatio": 0.0,
// 			"auctionExtension": 1
// 		  },
// 		  "logNormal": {
// 			"riskAversionParameter": 0.0001,
// 			"tau": 0.0000190129,
// 			"params": {
// 			  "mu": 0,
// 			  "r": 0.016,
// 			  "sigma": 1.25
// 			}
// 		  }
// 		}
// 	  }
// 	}
//   }