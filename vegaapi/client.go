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
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	vegaeventspb "code.vegaprotocol.io/vega/protos/vega/events/v1"
)

type VegaCoreClient interface {
	LastBlock(ctx context.Context) (*vegaapipb.LastBlockHeightResponse, error)
	Statistics(ctx context.Context) (*vegaapipb.StatisticsResponse, error)
	SendTransaction(ctx context.Context, tx *commandspb.Transaction, reqType vegaapipb.SubmitTransactionRequest_Type) (response *vegaapipb.SubmitTransactionResponse, err error)
	DepositBuiltinAsset(ctx context.Context, vegaAssetId string, partyId string, amount string, signAny func([]byte) ([]byte, string, error)) (bool, error)
	CoreNetworkParameters(ctx context.Context, parameterKey string) ([]*vega.NetworkParameter, error)
}

type DataNodeClient interface {
	VegaCoreClient
	ListNetworkParameters(ctx context.Context) (*types.NetworkParams, error)
	GetCurrentEpoch(ctx context.Context) (*vega.Epoch, error)
	ListAssets(ctx context.Context) (map[string]*vega.AssetDetails, error)
	ListMarkets(ctx context.Context) ([]*vega.Market, error)
	GetAllMarketsWithState(ctx context.Context, states []vega.Market_State) ([]*vega.Market, error)
	GetPartyTotalStake(ctx context.Context, partyId string) (*big.Int, error)
	GeneralAccountBalance(ctx context.Context, partyID, assetID string) (*big.Int, error)
	ListAccounts(ctx context.Context, partyID string, accountType vega.AccountType, assetId *string) ([]datanode.AccountFunds, error)
	ListCoreSnapshots(ctx context.Context) ([]vegaeventspb.CoreSnapshotData, error)
	LastNetworkHistorySegment(ctx context.Context) (*dataapipb.HistorySegment, error)
	ListProtocolUpgradeProposals(ctx context.Context) ([]vegaeventspb.ProtocolUpgradeEvent, error)
	ListGovernanceData(ctx context.Context, req *v2.ListGovernanceDataRequest) (response *v2.ListGovernanceDataResponse, err error)
	GetGovernanceData(ctx context.Context, req *v2.GetGovernanceDataRequest) (response *v2.GetGovernanceDataResponse, err error)
	ListVotes(ctx context.Context, req *v2.ListVotesRequest) (response *v2.ListVotesResponse, err error)
	GetCurrentReferralProgram(ctx context.Context) (*v2.ReferralProgram, error)
	ListReferralSets(ctx context.Context) (map[string]*v2.ReferralSet, error)
	ListReferralSetReferees(ctx context.Context) (map[string]v2.ReferralSetReferee, error)
	GetCurrentVolumeDiscountProgram(ctx context.Context) (*dataapipb.VolumeDiscountProgram, error)
}
