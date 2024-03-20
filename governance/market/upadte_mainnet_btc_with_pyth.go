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

func UpdateBTCUSDMainnetMarketProposal(
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	reference := tools.RandAlphaNumericString(40)

	contractABI := `[{
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

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "Update market BTC/USDT Perp Mainnet",
			Description: "## Summary\n\nThis proposal requests to list BTCUSDT Pyth Perp as a market with USDT as a settlement asset",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),

			Change: &vega.ProposalTerms_UpdateMarket{
				UpdateMarket: &vega.UpdateMarket{
					MarketId: "00788b763b999ef555ac5da17de155ff4237dd14aa6671a303d1285f27f094f0",
					Changes: &vega.UpdateMarketConfiguration{
						MarkPriceConfiguration: &vega.CompositePriceConfiguration{
							CompositePriceType:       vega.CompositePriceType_COMPOSITE_PRICE_TYPE_WEIGHTED,
							DecayWeight:              "1.0",
							DecayPower:               1,
							CashAmount:               "5000000",
							SourceWeights:            []string{"0.0", "0.0", "0.0", "1.0"},
							SourceStalenessTolerance: []string{"1m", "1m", "1m", "1m"},
							DataSourcesSpec: []*vega.DataSourceDefinition{
								{
									SourceType: &vega.DataSourceDefinition_External{
										External: &vega.DataSourceDefinitionExternal{
											SourceType: &vega.DataSourceDefinitionExternal_EthOracle{
												EthOracle: &vega.EthCallSpec{
													Address: "0x719abd606155442c21b7d561426d42bd0e40a776",
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
							},
							DataSourcesSpecBinding: []*vega.SpecBindingForCompositePrice{
								{
									PriceSourceProperty: "btc.price",
								},
							},
						},
						LinearSlippageFactor:    "0.001",
						QuadraticSlippageFactor: "0",
						Instrument: &vega.UpdateInstrumentConfiguration{
							Name: "BTC/USD(USDT)-Perp-Mainnet",
							Code: "BTC/USD-PERP-MAINNET",
							Product: &vega.UpdateInstrumentConfiguration_Perpetual{
								Perpetual: &vega.UpdatePerpetualProduct{
									InternalCompositePriceConfiguration: &vega.CompositePriceConfiguration{
										CompositePriceType:       vega.CompositePriceType_COMPOSITE_PRICE_TYPE_WEIGHTED,
										DecayWeight:              "1.0",
										DecayPower:               1,
										CashAmount:               "50000000",
										SourceWeights:            []string{"0.0", "1.0", "0.0"},
										SourceStalenessTolerance: []string{"1m", "1m", "1m"},
									},
									FundingRateScalingFactor: ptr.From("0.0625"),
									ClampLowerBound:          "-0.0005",
									ClampUpperBound:          "0.0005",
									InterestRate:             "0.1095",
									MarginFundingFactor:      "0.9",
									QuoteName:                "USDT",
									DataSourceSpecForSettlementData: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_External{
											External: &vega.DataSourceDefinitionExternal{
												SourceType: &vega.DataSourceDefinitionExternal_EthOracle{
													EthOracle: &vega.EthCallSpec{
														Address: "0x719abd606155442c21b7d561426d42bd0e40a776",
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
									DataSourceSpecBinding: &vega.DataSourceSpecToPerpetualBinding{
										SettlementDataProperty:     "btc.price",
										SettlementScheduleProperty: "vegaprotocol.builtin.timetrigger",
									},
								},
							},
						},
						Metadata: append([]string{
							"base:BTC",
							"quote:USD",
							"oracle:pyth",
							"chain:gnosis",
							"class:fx/crypto",
							"perpetual",
							"sector:defi",
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
						LiquiditySlaParameters: &vega.LiquiditySLAParameters{
							PriceRange:                  "0.03",
							CommitmentMinTimeFraction:   "0.685",
							PerformanceHysteresisEpochs: 1,
							SlaCompetitionFactor:        "0.5",
						},
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 0.05,
							},
							TriggeringRatio:  "0.9",
							AuctionExtension: 1,
						},
						LiquidationStrategy: &vega.LiquidationStrategy{
							DisposalTimeStep:    30,
							MaxFractionConsumed: "0.1",
							DisposalFraction:    "0.1",
							// FullDisposalSize:    0,
						},
						RiskParameters: &vega.UpdateMarketConfiguration_LogNormal{
							LogNormal: &vega.LogNormalRiskModel{
								RiskAversionParameter: 0.000001,
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
