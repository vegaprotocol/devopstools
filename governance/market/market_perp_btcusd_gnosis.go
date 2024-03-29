package market

import (
	"time"

	"code.vegaprotocol.io/vega/libs/ptr"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
	"github.com/vegaprotocol/devopstools/tools"
	"google.golang.org/protobuf/types/known/structpb"
)

const PerpetualBTCUSDGnosis = "auto:perpetual_btc_usd_gnosis"
const PerpetualBTCUSDOracleGnosisAddress = "0xf72475a8778E2e4b953f99a85f21a2feADbF77B9"

func NewBTCUSDGnosisPerpetualMarketProposal(
	settlementVegaAssetId string,
	oracleAddress string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	var (
		reference = tools.RandAlphaNumericString(40)
		name      = "BTCUSDT Gnosis Perp"
	)

	contractABI := `[{"inputs":[{"internalType":"string","name":"pair","type":"string"}],"name":"getPrice","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       name,
			Description: "## Summary\n\nThis proposal requests to list BTCUSDT Gnosis Perp as a market with USDT as a settlement asset",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						MarkPriceConfiguration: &vega.CompositePriceConfiguration{
							CompositePriceType:       vega.CompositePriceType_COMPOSITE_PRICE_TYPE_WEIGHTED,
							DecayWeight:              "1.0",
							DecayPower:               1,
							CashAmount:               "500000",
							SourceWeights:            []string{"0.999", "0.001", "0.0"},
							SourceStalenessTolerance: []string{"1m", "1m", "1m"},
						},
						DecimalPlaces:           4,
						PositionDecimalPlaces:   1,
						LinearSlippageFactor:    "0.001",
						QuadraticSlippageFactor: "0",
						TickSize:                "1",
						Instrument: &vega.InstrumentConfiguration{
							Name: name,
							Code: "BTCUSDT.GNOSIS.PERP",
							Product: &vega.InstrumentConfiguration_Perpetual{
								Perpetual: &vega.PerpetualProduct{
									InternalCompositePriceConfiguration: &vega.CompositePriceConfiguration{
										CompositePriceType:       vega.CompositePriceType_COMPOSITE_PRICE_TYPE_WEIGHTED,
										DecayWeight:              "1.0",
										DecayPower:               1,
										CashAmount:               "50000000",
										SourceWeights:            []string{"0.0", "1.0", "0.0"},
										SourceStalenessTolerance: []string{"1m", "1m", "1m"},
									},
									FundingRateScalingFactor: ptr.From("0.01041666667"),
									ClampLowerBound:          "-0.0005",
									ClampUpperBound:          "0.0005",
									InterestRate:             "0.1095",
									MarginFundingFactor:      "0.9",
									SettlementAsset:          settlementVegaAssetId,
									QuoteName:                "USDT",
									DataSourceSpecForSettlementData: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_External{
											External: &vega.DataSourceDefinitionExternal{
												SourceType: &vega.DataSourceDefinitionExternal_EthOracle{
													EthOracle: &vega.EthCallSpec{
														Address: oracleAddress,
														Abi:     contractABI,
														Method:  "getPrice",
														Args: []*structpb.Value{
															structpb.NewStringValue("BTCUSDT"),
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
							"chain:gnosis",
							"class:fx/crypto",
							"perpetual",
							"managed:vega/ops",
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
						LiquiditySlaParameters: &vega.LiquiditySLAParameters{
							PriceRange:                  "0.015",
							CommitmentMinTimeFraction:   "0.60",
							PerformanceHysteresisEpochs: 0,
							SlaCompetitionFactor:        "0.2",
						},
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 10,
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
						RiskParameters: &vega.NewMarketConfiguration_LogNormal{
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
