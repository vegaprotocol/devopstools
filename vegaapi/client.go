package vegaapi

import (
	"math/big"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	vegaeventspb "code.vegaprotocol.io/vega/protos/vega/events/v1"

	"github.com/vegaprotocol/devopstools/vegaapi/datanode"
)

type VegaCoreClient interface {
	LastBlockData() (*vegaapipb.LastBlockHeightResponse, error)
	Statistics() (*vegaapipb.StatisticsResponse, error)
	SubmitTransaction(
		req *vegaapipb.SubmitTransactionRequest,
	) (response *vegaapipb.SubmitTransactionResponse, err error)
	PropagateChainEvent(
		req *vegaapipb.PropagateChainEventRequest,
	) (response *vegaapipb.PropagateChainEventResponse, err error)
	DepositBuiltinAsset(
		vegaAssetId string,
		partyId string,
		amount string,
		signAny func([]byte) ([]byte, string, error),
	) (bool, error)
	DepositERC20Asset(
		vegaAssetId string,
		sourceEthereumAddress string,
		targetPartyId string,
		amount string,
		signAny func([]byte) ([]byte, string, error),
	) (bool, error)
	CoreNetworkParameters(parameterKey string) ([]*vega.NetworkParameter, error)
}

type DataNodeClient interface {
	VegaCoreClient
	GetAllNetworkParameters() (map[string]string, error)
	GetCurrentEpoch() (*vega.Epoch, error)
	GetAssets() (map[string]*vega.AssetDetails, error)
	GetAllMarkets() ([]*vega.Market, error)
	GetMarketById(marketId string) (*vega.Market, error)
	GetPartyTotalStake(partyId string) (*big.Int, error)
	GetFunds(
		partyID string,
		accountType vega.AccountType,
		assetId *string,
	) ([]datanode.AccountFunds, error)
	ListCoreSnapshots() ([]vegaeventspb.CoreSnapshotData, error)
	LastNetworkHistorySegment() (*dataapipb.HistorySegment, error)
	ListProtocolUpgradeProposals() ([]vegaeventspb.ProtocolUpgradeEvent, error)
	ListGovernanceData(
		req *dataapipb.ListGovernanceDataRequest,
	) (response *dataapipb.ListGovernanceDataResponse, err error)
}
