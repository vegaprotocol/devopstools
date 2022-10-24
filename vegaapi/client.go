package vegaapi

import (
	"math/big"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
)

type VegaClient interface {
	LastBlockData() (*vegaapipb.LastBlockHeightResponse, error)
	Statistics() (*vegaapipb.StatisticsResponse, error)
	SubmitTransaction(req *vegaapipb.SubmitTransactionRequest) (response *vegaapipb.SubmitTransactionResponse, err error)
}

type DataNodeClient interface {
	VegaClient
	GetAllNetworkParameters() (map[string]string, error)
	GetCurrentEpoch() (*vega.Epoch, error)
	GetAssets() (map[string]*vega.AssetDetails, error)
	GetAllMarkets() ([]*vega.Market, error)
	GetPartyTotalStake(partyId string) (*big.Int, error)

	ListGovernanceData(req *dataapipb.ListGovernanceDataRequest) (response *dataapipb.ListGovernanceDataResponse, err error)
}
