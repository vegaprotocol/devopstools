package vegaapi

import (
	"math/big"

	"code.vegaprotocol.io/vega/protos/vega"
)

type MarketInfo struct {
	Market          *vega.Market
	MarketData      *vega.MarketData
	SettlementAsset *vega.Asset
}

type AccountFunds struct {
	Balance     *big.Int
	PartyId     string
	AccountType vega.AccountType
	AssetId     string
}
