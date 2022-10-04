package proposals

import (
	"fmt"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	oraclespb "code.vegaprotocol.io/vega/protos/vega/oracles/v1"
	"github.com/vegaprotocol/devopstools/tools"
)

func NewTSLAMarketProposal(
	settlementVegaAssetId string,
	decimalPlaces uint64,
	oraclePubKey string,
	closingTime time.Time,
	enactmentTime time.Time,
	extraMetadata []string,
) *commandspb.ProposalSubmission {
	var (
		reference = tools.RandAlpaNumericString(40)
		Name      = fmt.Sprintf("Tesla Quarterly (%s)", time.Now().AddDate(0, 3, 0).Format("Jan 2006")) // Now + 3 months
	)

	return &commandspb.ProposalSubmission{
		Reference: reference,
		Rationale: &vega.ProposalRationale{
			Title:       "New EURO market",
			Description: "New EURO market",
		},
		Terms: &vega.ProposalTerms{
			ClosingTimestamp:   closingTime.Unix(),
			EnactmentTimestamp: enactmentTime.Unix(),
			Change: &vega.ProposalTerms_NewMarket{
				NewMarket: &vega.NewMarket{
					Changes: &vega.NewMarketConfiguration{
						DecimalPlaces: decimalPlaces,
						//PositionDecimalPlaces: 3,
						Instrument: &vega.InstrumentConfiguration{
							Name: Name,
							Code: "TSLA.QM21",
							Product: &vega.InstrumentConfiguration_Future{
								Future: &vega.FutureProduct{
									SettlementAsset: settlementVegaAssetId,
									QuoteName:       "EURO",
									OracleSpecForSettlementPrice: &oraclespb.OracleSpecConfiguration{
										PubKeys: []string{oraclePubKey},
										Filters: []*oraclespb.Filter{
											{
												Key: &oraclespb.PropertyKey{
													Name: "prices.TSLA.value",
													Type: oraclespb.PropertyKey_TYPE_INTEGER,
												},
												Conditions: []*oraclespb.Condition{
													{
														Operator: oraclespb.Condition_OPERATOR_EQUALS,
														Value:    "1",
													},
												},
											},
										},
									},
									OracleSpecForTradingTermination: &oraclespb.OracleSpecConfiguration{
										PubKeys: []string{oraclePubKey},
										Filters: []*oraclespb.Filter{
											{
												Key: &oraclespb.PropertyKey{
													Name: "termination.TSLA.value",
													Type: oraclespb.PropertyKey_TYPE_BOOLEAN,
												},
												Conditions: []*oraclespb.Condition{
													{
														Operator: oraclespb.Condition_OPERATOR_EQUALS,
														Value:    "1",
													},
												},
											},
										},
									},
									OracleSpecBinding: &vega.OracleSpecToFutureBinding{
										SettlementPriceProperty:    "prices.TSLA.value",
										TradingTerminationProperty: "termination.TSLA.value",
									},
								},
							},
						},
						Metadata: append([]string{
							"formerly:5A86B190C384997F",
							"quote:EURO",
							"ticker:TSLA",
							"class:equities/single-stock-futures",
							"sector:tech",
							"listing_venue:NASDAQ",
							"country:US",
						}, extraMetadata...),
						PriceMonitoringParameters: &vega.PriceMonitoringParameters{
							Triggers: []*vega.PriceMonitoringTrigger{
								{
									Horizon:          43200,
									Probability:      "0.9999999",
									AuctionExtension: 600,
								},
							},
						},
						LiquidityMonitoringParameters: &vega.LiquidityMonitoringParameters{
							TargetStakeParameters: &vega.TargetStakeParameters{
								TimeWindow:    3600,
								ScalingFactor: 10,
							},
							TriggeringRatio:  0.7,
							AuctionExtension: 1,
						},
						RiskParameters: &vega.NewMarketConfiguration_LogNormal{
							LogNormal: &vega.LogNormalRiskModel{
								RiskAversionParameter: 0.01,
								Tau:                   0.0001140771161,
								Params: &vega.LogNormalModelParams{
									Mu:    0,
									R:     0.016,
									Sigma: 0.8,
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
