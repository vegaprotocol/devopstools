package proposals

import (
	"time"

	"code.vegaprotocol.io/vega/libs/ptr"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

const PerpetualLINKUSD = "auto:perpetual_link_usd"
const PerpetualLINKUSDOracleAddress = "0xc59E3633BAAC79493d908e63626716e204A45EdF"

func NewLINKUSDPerpetualMarketProposal(
	settlementVegaAssetId string,
	decimalPlaces uint64,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	var (
		reference = tools.RandAlpaNumericString(40)
		name      = "LINKUSD Perpetual"
	)

	contractABI := `[{"inputs":[],"name":"latestAnswer","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"}]`

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "New LINKUSD perpetual market",
			Description: "New LINKUSD perpetual market",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						DecimalPlaces:           decimalPlaces,
						PositionDecimalPlaces:   1,
						LinearSlippageFactor:    "0.01",
						QuadraticSlippageFactor: "0.0",
						Instrument: &vega.InstrumentConfiguration{
							Name: name,
							Code: "LINKUSD.PERP",
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
														Address: oraclePubKey, // chainlink LINK/USD
														Abi:     contractABI,
														Method:  "latestAnswer",
														Normalisers: []*vega.Normaliser{
															{
																Name:       "link.price",
																Expression: "$[0]",
															},
														},
														RequiredConfirmations: 3,
														Trigger: &vega.EthCallTrigger{
															Trigger: &vega.EthCallTrigger_TimeTrigger{
																TimeTrigger: &vega.EthTimeTrigger{
																	Every: ptr.From(uint64(120)), // every 2 mins
																},
															},
														},
														Filters: []*datav1.Filter{
															{
																Key: &datav1.PropertyKey{
																	Name:                "link.price",
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
																Every: 14400, // 4hrs in seconds
															},
														},
													},
												},
											},
										},
									},
									DataSourceSpecBinding: &vega.DataSourceSpecToPerpetualBinding{
										SettlementDataProperty:     "link.price",
										SettlementScheduleProperty: "vegaprotocol.builtin.timetrigger",
									},
								},
							},
						},
						Metadata: append([]string{
							"formerly:706572706c696e6b757364",
							"base:LINK",
							"quote:USD",
							"class:fx/crypto",
							"monthly",
							"sector:crypto",
						}, extraMetadata...),
						PriceMonitoringParameters: &vega.PriceMonitoringParameters{
							Triggers: []*vega.PriceMonitoringTrigger{
								{
									Horizon:          4320,
									Probability:      "0.99",
									AuctionExtension: 300,
								},
								{
									Horizon:          1440,
									Probability:      "0.99",
									AuctionExtension: 180,
								},
								{
									Horizon:          360,
									Probability:      "0.99",
									AuctionExtension: 120,
								},
							},
						},
						LpPriceRange: "0.05",
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 10,
							},
							TriggeringRatio:  "0.9",
							AuctionExtension: 1,
						},
						RiskParameters: &vega.NewMarketConfiguration_LogNormal{
							LogNormal: &vega.LogNormalRiskModel{
								RiskAversionParameter: 0.0001,
								Tau:                   0.00000380258,
								Params: &vega.LogNormalModelParams{
									Mu:    0,
									R:     0,
									Sigma: 1.5,
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
// 	  "title": "LINKUSD market",
// 	  "description": "New LINKUSD market"
// 	},
// 	"terms": {
// 	  "closingTimestamp": 0,
// 	  "enactmentTimestamp": 0,
// 	  "newMarket": {
// 		"changes": {
// 		  "instrument": {
// 			"name": "LINKUSD Monthly (30 Jun 2022)",
// 			"code": "LINKUSD.MF21",
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
// 					  "name": "prices.LINK.value",
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
// 					  "name": "termination.LINK.value",
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
// 				"settlementPriceProperty": "prices.LINK.value",
// 				"tradingTerminationProperty": "termination.LINK.value"
// 			  }
// 			}
// 		  },
// 		  "decimalPlaces": 5,
// 		  "positionDecimalPlaces": 3,
// 		  "metadata": [
// 			"formerly:076BB86A5AA41E3E",
// 			"base:LINK",
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
