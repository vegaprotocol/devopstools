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

func NewMainnetSimulationLDOUSDTPerpWithoutTime(settlementAsset string) *commandspb.ProposalSubmission {
	return NewMainnetSimulationLDOUSDTPerp(
		settlementAsset,
		time.Now().Add(30*time.Second),
		time.Now().Add(45*time.Second),
		DefaultExtraMetadata,
	)
}

func NewMainnetSimulationLDOUSDTPerp(
	settlementVegaAssetId string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	var (
		reference = tools.RandAlphaNumericString(40)
	)
	const (
		contractAddress = "0x9eb2EBD260D82410592758B3091F74977E4A404c"
		contractABI     = `[{
			"inputs": [
				{
					"internalType": "contract IUniswapV3Pool",
					"name": "pool", 
					"type": "address"
				},
				{
					"internalType": "uint32",
					"name": "twapInterval",
					"type": "uint32"
				}
			],
			"name": "priceFromEthPoolInUsdt",
			"outputs": [
				{
					"internalType": "uint256", 
					"name": "",
					"type": "uint256"
				}
			],
			"stateMutability": "view",
			"type":"function"
		}]`
	)

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "VMP-26 Create market LDO/USDT Perpetual",
			Description: "## Summary\n\nThis proposal requests to list LDO/USDT Perpetual as a market with USDT as a settlement asset",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						DecimalPlaces:           3,
						PositionDecimalPlaces:   1,
						LinearSlippageFactor:    "0.001",
						QuadraticSlippageFactor: "0.0",
						TickSize:                "1",
						Instrument: &vega.InstrumentConfiguration{
							Name: "LDO/USDT-Perp",
							Code: "LDO/USDT-PERP",
							Product: &vega.InstrumentConfiguration_Perpetual{
								Perpetual: &vega.PerpetualProduct{
									SettlementAsset:     settlementVegaAssetId,
									QuoteName:           "USDT",
									MarginFundingFactor: "0.9",
									ClampLowerBound:     "-0.0005",
									ClampUpperBound:     "0.0005",
									InterestRate:        "0.1095",
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
																Every: 28800, // 24h
															},
														},
													},
												},
											},
										},
									},
									DataSourceSpecBinding: &vega.DataSourceSpecToPerpetualBinding{
										SettlementDataProperty:     "ldo.price",
										SettlementScheduleProperty: "vegaprotocol.builtin.timetrigger",
									},
									DataSourceSpecForSettlementData: &vega.DataSourceDefinition{
										SourceType: &vega.DataSourceDefinition_External{
											External: &vega.DataSourceDefinitionExternal{
												SourceType: &vega.DataSourceDefinitionExternal_EthOracle{
													EthOracle: &vega.EthCallSpec{
														// https://docs.chain.link/data-feeds/price-feeds/addresses#Sepolia%20Testnet
														Address: contractAddress,
														Abi:     contractABI,
														Method:  "priceFromEthPoolInUsdt",
														Args: []*structpb.Value{
															structpb.NewStringValue("0xa3f558aebaecaf0e11ca4b2199cc5ed341edfd74"),
															structpb.NewNumberValue(300),
														},
														Normalisers: []*vega.Normaliser{
															{
																Name:       "ldo.price",
																Expression: "$[0]",
															},
														},
														RequiredConfirmations: 3,
														Trigger: &vega.EthCallTrigger{
															Trigger: &vega.EthCallTrigger_TimeTrigger{
																TimeTrigger: &vega.EthTimeTrigger{
																	Every: ptr.From(uint64(300)),
																},
															},
														},
														Filters: []*datav1.Filter{
															{
																Key: &datav1.PropertyKey{
																	Name:                "ldo.price",
																	Type:                datav1.PropertyKey_TYPE_INTEGER,
																	NumberDecimalPlaces: ptr.From(uint64(18)),
																},
																Conditions: []*datav1.Condition{
																	{
																		Operator: datav1.Condition_OPERATOR_GREATER_THAN_OR_EQUAL,
																		Value:    "0",
																	},
																},
															},
														},
														SourceChainId: SepoliaChainID,
													},
												},
											},
										},
									},
								},
							},
						},
						Metadata: append([]string{
							"base:LDO",
							"quote:USDT",
							"class:fx/crypto",
							"perpetual",
							"sector:defi",
							"enactment:2024-02-03T14:00:00Z",
						}, extraMetadata...),
						PriceMonitoringParameters: &vega.PriceMonitoringParameters{
							Triggers: []*vega.PriceMonitoringTrigger{
								{
									Horizon:          43200,
									Probability:      "0.9999999",
									AuctionExtension: 300,
								},
								{
									Horizon:          1440,
									Probability:      "0.9999999",
									AuctionExtension: 180,
								},
								{
									Horizon:          360,
									Probability:      "0.9999999",
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
								ScalingFactor: 0.05,
							},
							TriggeringRatio:  "0.9",
							AuctionExtension: 1,
						},
						LiquidationStrategy: &vega.LiquidationStrategy{
							DisposalTimeStep:    30,
							MaxFractionConsumed: "0.1",
							DisposalFraction:    "0.1",
							FullDisposalSize:    0,
						},
						RiskParameters: &vega.NewMarketConfiguration_LogNormal{
							LogNormal: &vega.LogNormalRiskModel{
								RiskAversionParameter: 0.000001,
								Tau:                   0.0000071,
								Params: &vega.LogNormalModelParams{
									Mu:    0,
									R:     0,
									Sigma: 1.5,
								},
							},
						},
						MarkPriceConfiguration: &vega.CompositePriceConfiguration{
							CompositePriceType: vega.CompositePriceType_COMPOSITE_PRICE_TYPE_LAST_TRADE,
						},
					},
				},
			},
		},
	}
}
