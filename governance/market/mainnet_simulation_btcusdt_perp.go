package market

import (
	"time"

	"github.com/vegaprotocol/devopstools/tools"

	"code.vegaprotocol.io/vega/libs/ptr"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"

	"google.golang.org/protobuf/types/known/structpb"
)

func NewMainnetSimulationBitcoinTetherPerpetualWithoutTime(settlementAsset string) *commandspb.ProposalSubmission {
	return NewMainnetSimulationBitcoinTetherPerpetual(
		settlementAsset,
		time.Now().Add(30*time.Second),
		time.Now().Add(45*time.Second),
	)
}

func NewMainnetSimulationBitcoinTetherPerpetual(
	settlementAsset string,
	closingTime time.Time,
	enactmentTime time.Time,
) *commandspb.ProposalSubmission {
	reference := tools.RandAlphaNumericString(40)
	const (
		contractAddress = "0x719abd606155442c21b7d561426d42bd0e40a776"
		contractABI     = `[{
			"inputs": [
			  {
				"internalType": "bytes32",
				"name": "id",
				"type": "bytes32"
			  }
			],
			"name": "getPrice",
			"outputs": [
			  {
				"internalType": "int256",
				"name": "",
				"type": "int256"
			  }
			],
			"stateMutability": "view",
			"type": "function"
			}]`
	)

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						TickSize:              "1",
						PositionDecimalPlaces: 4,
						DecimalPlaces:         1,
						Instrument: &vega.InstrumentConfiguration{
							Name: "Bitcoin / Tether USD (Perpetual)",
							Code: "BTCUSDT.PERP",
							Product: &vega.InstrumentConfiguration_Perpetual{
								Perpetual: &vega.PerpetualProduct{
									SettlementAsset:          settlementAsset,
									QuoteName:                "USDT",
									MarginFundingFactor:      "0.9",
									InterestRate:             "0.1095",
									ClampLowerBound:          "-0.0005",
									ClampUpperBound:          "0.0005",
									FundingRateScalingFactor: ptr.From("1"),
									FundingRateLowerBound:    ptr.From("-0.01"),
									FundingRateUpperBound:    ptr.From("0.01"),
									DataSourceSpecForSettlementSchedule: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_Internal{
											Internal: &vega.DataSourceDefinitionInternal{
												SourceType: &vega.DataSourceDefinitionInternal_TimeTrigger{
													TimeTrigger: &vega.DataSourceSpecConfigurationTimeTrigger{
														Conditions: []*datav1.Condition{
															{
																Operator: datav1.Condition_OPERATOR_GREATER_THAN,
																Value:    "0",
															},
														},
														Triggers: []*datav1.InternalTimeTrigger{
															{
																Every: 28800, // 8hrs
															},
														},
													},
												},
											},
										},
									},
									DataSourceSpecForSettlementData: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_External{
											External: &vega.DataSourceDefinitionExternal{
												SourceType: &vega.DataSourceDefinitionExternal_EthOracle{
													EthOracle: &vega.EthCallSpec{
														Address: contractAddress,
														Abi:     contractABI,
														Method:  "getPrice",
														Args: []*structpb.Value{
															structpb.NewStringValue("0xe62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43"),
														},
														Normalisers: []*vega.Normaliser{
															{
																Name:       "btc.price",
																Expression: "$[0]",
															},
														},
														RequiredConfirmations: 3,
														Trigger: &vega.EthCallTrigger{
															Trigger: &vega.EthCallTrigger_TimeTrigger{
																TimeTrigger: &vega.EthTimeTrigger{
																	Every: ptr.From(uint64(60)), // pool prices every minutes
																},
															},
														},
														Filters: []*datav1.Filter{
															{
																Key: &datav1.PropertyKey{
																	Name:                "btc.price",
																	Type:                datav1.PropertyKey_TYPE_INTEGER,
																	NumberDecimalPlaces: ptr.From(uint64(18)),
																},
																Conditions: []*datav1.Condition{
																	{
																		Operator: datav1.Condition_OPERATOR_GREATER_THAN,
																		Value:    "0",
																	},
																},
															},
														},
														SourceChainId: GnosisChainID,
													},
												},
											},
										},
									},
									DataSourceSpecBinding: &vega.DataSourceSpecToPerpetualBinding{
										SettlementDataProperty:     "btc.price",
										SettlementScheduleProperty: "vegaprotocol.builtin.timetrigger",
									},
									InternalCompositePriceConfiguration: &vega.CompositePriceConfiguration{
										DecayWeight: "1",
										DecayPower:  1,
										CashAmount:  "50000000",
										SourceWeights: []string{
											"0", "0.999", "0.001", "0",
										},
										SourceStalenessTolerance: []string{
											"1m0s", "1m0s", "10m0s", "10m0s",
										},
										CompositePriceType: vega.CompositePriceType_COMPOSITE_PRICE_TYPE_WEIGHTED,
										DataSourcesSpec: []*vega.DataSourceDefinition{
											{
												SourceType: &vega.DataSourceDefinition_External{
													External: &vega.DataSourceDefinitionExternal{
														SourceType: &vega.DataSourceDefinitionExternal_EthOracle{
															EthOracle: &vega.EthCallSpec{
																Address: contractAddress,
																Abi:     contractABI,
																Method:  "getPrice",
																Args: []*structpb.Value{
																	structpb.NewStringValue("0xe62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43"),
																},
																Trigger: &vega.EthCallTrigger{
																	Trigger: &vega.EthCallTrigger_TimeTrigger{
																		TimeTrigger: &vega.EthTimeTrigger{
																			Every: ptr.From(uint64(60)), // pool prices every minutes
																		},
																	},
																},
																RequiredConfirmations: 3,
																Filters: []*datav1.Filter{
																	{
																		Key: &datav1.PropertyKey{
																			Name:                "btc.price",
																			Type:                datav1.PropertyKey_TYPE_INTEGER,
																			NumberDecimalPlaces: ptr.From(uint64(18)),
																		},
																		Conditions: []*datav1.Condition{
																			{
																				Operator: datav1.Condition_OPERATOR_GREATER_THAN,
																				Value:    "0",
																			},
																		},
																	},
																},
																Normalisers: []*vega.Normaliser{
																	{
																		Name:       "btc.price",
																		Expression: "$[0]",
																	},
																},
																SourceChainId: GnosisChainID,
															},
														},
													},
												},
											},
										},
										DataSourcesSpecBinding: []*vega.SpecBindingForCompositePrice{
											{
												PriceSourceProperty: "btc.price",
											},
										},
									},
								},
							},
						},
						Metadata: []string{
							"base:BTC",
							"quote:USD",
							"oracle:pyth",
							"chain:gnosis",
							"class:fx/crypto",
							"perpetual",
							"sector:defi",
						},
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
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 0.05,
							},
							TriggeringRatio:  "0.9",
							AuctionExtension: 1,
						},
						RiskParameters: &vega.NewMarketConfiguration_LogNormal{
							LogNormal: &vega.LogNormalRiskModel{
								RiskAversionParameter: 1e-06,
								Tau:                   3.80258e-06,
								Params: &vega.LogNormalModelParams{
									Mu:    0,
									R:     0,
									Sigma: 1.5,
								},
							},
						},
						LinearSlippageFactor:    "0.001",
						QuadraticSlippageFactor: "0",
						LiquiditySlaParameters: &vega.LiquiditySLAParameters{
							PriceRange:                  "0.03",
							CommitmentMinTimeFraction:   "0.685",
							PerformanceHysteresisEpochs: 1,
							SlaCompetitionFactor:        "0.5",
						},
						LiquidationStrategy: &vega.LiquidationStrategy{
							DisposalTimeStep:      30,
							DisposalFraction:      "0.1",
							FullDisposalSize:      0,
							MaxFractionConsumed:   "0.1",
							DisposalSlippageRange: "0.005",
						},
						MarkPriceConfiguration: &vega.CompositePriceConfiguration{
							DecayWeight: "1.0",
							DecayPower:  1,
							CashAmount:  "5000000",
							SourceWeights: []string{
								"0.0", "0.0", "0.0", "1.0",
							},
							SourceStalenessTolerance: []string{
								"1m", "1m", "1m", "1m",
							},
							CompositePriceType: vega.CompositePriceType_COMPOSITE_PRICE_TYPE_WEIGHTED,
							DataSourcesSpec: []*vega.DataSourceDefinition{
								{
									SourceType: &vega.DataSourceDefinition_External{
										External: &vega.DataSourceDefinitionExternal{
											SourceType: &vega.DataSourceDefinitionExternal_EthOracle{
												EthOracle: &vega.EthCallSpec{
													Address: contractAddress,
													Abi:     contractABI,
													Method:  "getPrice",
													Args: []*structpb.Value{
														structpb.NewStringValue("0xe62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43"),
													},
													Trigger: &vega.EthCallTrigger{
														Trigger: &vega.EthCallTrigger_TimeTrigger{
															TimeTrigger: &vega.EthTimeTrigger{
																Every: ptr.From(uint64(60)), // pool prices every minutes
															},
														},
													},
													RequiredConfirmations: 3,
													Filters: []*datav1.Filter{
														{
															Key: &datav1.PropertyKey{
																Name:                "btc.price",
																Type:                datav1.PropertyKey_TYPE_INTEGER,
																NumberDecimalPlaces: ptr.From(uint64(18)),
															},
															Conditions: []*datav1.Condition{
																{
																	Operator: datav1.Condition_OPERATOR_GREATER_THAN,
																	Value:    "0",
																},
															},
														},
													},
													Normalisers: []*vega.Normaliser{
														{
															Name:       "btc.price",
															Expression: "$[0]",
														},
													},
													SourceChainId: GnosisChainID,
												},
											},
										},
									},
								},
							},
							DataSourcesSpecBinding: []*vega.SpecBindingForCompositePrice{
								{
									PriceSourceProperty: "btc.price",
								},
							},
						},
					},
				},
			},
		},
	}
}
