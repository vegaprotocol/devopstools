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

const PerpetualSNXUSD = "auto:perpetual_snx_usd"

func NewSNXUSDPerpetualMarketProposal(
	settlementVegaAssetId string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	var (
		reference = tools.RandAlphaNumericString(40)
		name      = "SNXUSD Perpetual"
	)

	contractABI := `[{"inputs":[{"internalType":"contract IUniswapV3Pool","name":"pool","type":"address"},{"internalType":"uint32","name":"twapInterval","type":"uint32"}],"name":"priceFromEthPoolInUsdt","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "New SNXUSD perpetual market",
			Description: "New SNXUSD perpetual market",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						MarkPriceConfiguration: &vega.CompositePriceConfiguration{
							CompositePriceType: vega.CompositePriceType_COMPOSITE_PRICE_TYPE_LAST_TRADE,
						},
						DecimalPlaces:           3,
						PositionDecimalPlaces:   1,
						LinearSlippageFactor:    "0.001",
						QuadraticSlippageFactor: "0.0",
						Instrument: &vega.InstrumentConfiguration{
							Name: name,
							Code: "SNXUSD.PERP",
							Product: &vega.InstrumentConfiguration_Perpetual{
								Perpetual: &vega.PerpetualProduct{
									ClampLowerBound:     "-0.0005",
									ClampUpperBound:     "0.0005",
									InterestRate:        "0.1095",
									MarginFundingFactor: "0.9",
									SettlementAsset:     settlementVegaAssetId,
									QuoteName:           "USD",
									DataSourceSpecForSettlementData: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_External{
											External: &vega.DataSourceDefinitionExternal{
												SourceType: &vega.DataSourceDefinitionExternal_EthOracle{
													EthOracle: &vega.EthCallSpec{
														Address: "0x9eb2ebd260d82410592758b3091f74977e4a404c", // Uniswap price source contrat
														Abi:     contractABI,
														Method:  "priceFromEthPoolInUsdt",
														Args: []*structpb.Value{
															structpb.NewStringValue(
																"0xede8dd046586d22625ae7ff2708f879ef7bdb8cf", // SNX/USDC pool
															),
															structpb.NewNumberValue(float64(300)), // 300 seconds twap
														},
														Normalisers: []*vega.Normaliser{
															{
																Name:       "snx.price",
																Expression: "$[0]",
															},
														},
														RequiredConfirmations: 3,
														Trigger: &vega.EthCallTrigger{
															Trigger: &vega.EthCallTrigger_TimeTrigger{
																TimeTrigger: &vega.EthTimeTrigger{
																	Every: ptr.From(uint64(60)), // every 1 mins
																},
															},
														},
														Filters: []*datav1.Filter{
															{
																Key: &datav1.PropertyKey{
																	Name:                "snx.price",
																	Type:                datav1.PropertyKey_TYPE_INTEGER,
																	NumberDecimalPlaces: ptr.From(uint64(18)), // All prices are 18 decimals on the uniswap contract
																},
																Conditions: []*datav1.Condition{
																	{
																		Operator: datav1.Condition_OPERATOR_GREATER_THAN_OR_EQUAL,
																		Value:    "0",
																	},
																},
															},
														},
														SourceChainId: EthereumMainnetID,
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
																Every: 28800, // 8hrs fundiong period in seconds
															},
														},
													},
												},
											},
										},
									},
									DataSourceSpecBinding: &vega.DataSourceSpecToPerpetualBinding{
										SettlementDataProperty:     "snx.price",
										SettlementScheduleProperty: "vegaprotocol.builtin.timetrigger",
									},
								},
							},
						},
						Metadata: append([]string{
							"formerly:706572706c696e6b757364",
							"base:SNX",
							"quote:USD",
							"class:fx/crypto",
							"monthly",
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
							PriceRange:                  "0.03",
							CommitmentMinTimeFraction:   "0.85",
							PerformanceHysteresisEpochs: 1,
							SlaCompetitionFactor:        "0.5",
						},
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 0.5,
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
								RiskAversionParameter: 0.0001,
								Tau:                   0.0000071,
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
