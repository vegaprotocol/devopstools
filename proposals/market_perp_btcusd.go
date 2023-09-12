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
		name      = "BTCUSD Perpetual Futures"
	)

	contractABI := `[{"inputs":[],"name":"latestAnswer","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"}]`

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "New BTCUSD perpetual futures market",
			Description: "New BTCUSD perpetual futures market",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						DecimalPlaces:           decimalPlaces,
						PositionDecimalPlaces:   4,
						LinearSlippageFactor:    "0.01",
						QuadraticSlippageFactor: "0.0",
						Instrument: &vega.InstrumentConfiguration{
							Name: name,
							Code: "BTCUSD.PERP",
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
														RequiredConfirmations: 3,
														Trigger: &vega.EthCallTrigger{
															Trigger: &vega.EthCallTrigger_TimeTrigger{
																TimeTrigger: &vega.EthTimeTrigger{
																	Every: ptr.From(uint64(30)),
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
																		Operator: datav1.Condition_OPERATOR_GREATER_THAN,
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
							"class:fx/crypto",
							"perpetual",
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
							TriggeringRatio:  "0.9",
							AuctionExtension: 1,
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
