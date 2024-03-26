package vegaapi

import (
	"context"
	"math/big"

	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	vegaeventspb "code.vegaprotocol.io/vega/protos/vega/events/v1"
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
	GetAllNetworkParameters() (*types.NetworkParams, error)
	ListNetworkParameters(req *dataapipb.ListNetworkParametersRequest) (response *dataapipb.ListNetworkParametersResponse, err error)
	GetCurrentEpoch() (*vega.Epoch, error)
	ListAssets(ctx context.Context) (map[string]*vega.AssetDetails, error)
	GetAllMarkets() ([]*vega.Market, error)
	GetMarketById(marketId string) (*vega.Market, error)
	GetPartyTotalStake(partyId string) (*big.Int, error)
	ListAccounts(ctx context.Context, partyID string, accountType vega.AccountType, assetId *string) ([]datanode.AccountFunds, error)
	ListCoreSnapshots() ([]vegaeventspb.CoreSnapshotData, error)
	LastNetworkHistorySegment() (*dataapipb.HistorySegment, error)
	ListProtocolUpgradeProposals() ([]vegaeventspb.ProtocolUpgradeEvent, error)
	ListGovernanceData(
		req *dataapipb.ListGovernanceDataRequest,
	) (response *dataapipb.ListGovernanceDataResponse, err error)
	GetGovernanceData(req *dataapipb.GetGovernanceDataRequest) (response *dataapipb.GetGovernanceDataResponse, err error)
	ListVotes(req *dataapipb.ListVotesRequest) (response *dataapipb.ListVotesResponse, err error)
	GetCurrentReferralProgram() (*dataapipb.ReferralProgram, error)
	GetReferralSets() (map[string]*v2.ReferralSet, error)
	GetReferralSetReferees() (map[string]v2.ReferralSetReferee, error)
	GetCurrentVolumeDiscountProgram() (*dataapipb.VolumeDiscountProgram, error)
}
