package market

import (
	"fmt"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	datav1 "code.vegaprotocol.io/vega/protos/vega/data/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

const MainnetLinkUSDTMetadataID = "auto:mainnet_link_usdt"

var (
	MainnetLinkUSDTOracleDecimalPlaces   = uint64(6)
	MainnetLinkUSDTMarketSettlementHours = 24 * 31
	MainnetLinkUSDTPositionDecimalPlaces = int64(1)
	MainnetLinkUSDTDecimalPlaces         = uint64(3)
)

// settlementVegaAssetId ideally with 6 decimal places
func NewMainnetLinkUSDT(
	settlementVegaAssetId string,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) (*commandspb.ProposalSubmission, error) {

	reference := tools.RandAlphaNumericString(40)

	nowTime := time.Now()
	settlementTime := nowTime.Add(time.Hour * (time.Duration(MainnetLinkUSDTMarketSettlementHours)))

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title: "LINK/USDT",
			Description: `# Summary

This proposal requests to list LINK/USDT-231231 as a market with USDT as a settlement asset on the Vega Network as discussed in: https://community.vega.xyz/.

# Rationale

	* LINK is among the top 20 largest Crypto asset with the highest volume and Marketcap.
	* Given the price, 3 decimal places will be used for price due to the number of valid digits in asset price.
	* Position decimal places will be set to 2 considering the value per contract
	* USDT is chosen as settlement asset due to its stability.`,
		},

		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						Instrument: &vega.InstrumentConfiguration{
							Name: "LINK/USDT",
							Code: "LINK/USDT-MAINNET",
							Product: &vega.InstrumentConfiguration_Future{
								Future: &vega.FutureProduct{
									SettlementAsset: settlementVegaAssetId,
									QuoteName:       "USDC",
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
																	Name:                "prices.LINK.value",
																	Type:                datav1.PropertyKey_TYPE_INTEGER,
																	NumberDecimalPlaces: &MainnetLinkUSDTOracleDecimalPlaces,
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
																	Name: "prices.LINK.timestamp",
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
										SettlementDataProperty:     "prices.LINK.value",
										TradingTerminationProperty: "vegaprotocol.builtin.timestamp",
									},
								},
							},
						},
						DecimalPlaces: MainnetLinkUSDTDecimalPlaces,
						Metadata: append([]string{
							"base:LINK",
							"quote:USDT",
							"class:fx/crypto",
							"quarterly",
							"sector:defi",
							fmt.Sprintf("enactment:%s", enactmentTime.Format(time.RFC3339)),
							MainnetLinkUSDTMetadataID,
							fmt.Sprintf("settlement:%s   2023-12-31T08:00:00Z", settlementTime.Format(time.RFC3339)),
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
						PositionDecimalPlaces:   MainnetLinkUSDTPositionDecimalPlaces,
						LpPriceRange:            &[]string{"0.05"}[0],
						LiquiditySlaParameters:  nil,
						LinearSlippageFactor:    "0.001",
						QuadraticSlippageFactor: "0",
					},
				},
			},
		},
	}, nil
}
