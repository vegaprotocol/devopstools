package datanode

import (
	"context"
	"fmt"
	"math/big"

	e "github.com/vegaprotocol/devopstools/errors"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"

	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetPartyTotalStake(partyId string) (*big.Int, error) {
	res, err := n.GetStake(&dataapipb.GetStakeRequest{
		PartyId: partyId,
	})
	if err != nil {
		return nil, err
	}
	totalStake := new(big.Int)
	totalStake, ok := totalStake.SetString(res.CurrentStakeAvailable, 0)
	if !ok {
		return nil, fmt.Errorf("failed to convert %s to big.Int", res.CurrentStakeAvailable)
	}

	return totalStake, nil
}

func (n *DataNode) GetPartyDelegationToNode(partyId string, nodeId string) (*big.Int, error) {
	epoch, err := n.GetCurrentEpoch()
	if err != nil {
		return nil, err
	}
	amount := "0"
	for _, delegation := range epoch.Delegations {
		if delegation.Party == partyId && delegation.NodeId == nodeId {
			amount = delegation.Amount
			break
		}
	}
	result := new(big.Int)
	result, ok := result.SetString(amount, 0)
	if !ok {
		return nil, fmt.Errorf("failed to convert %s to big.Int", amount)
	}
	return result, nil
}

// GetStake returns stakes for the given party.
func (n *DataNode) GetStake(req *dataapipb.GetStakeRequest) (response *dataapipb.GetStakeResponse, err error) {
	msg := "gRPC call failed (data-node): GetStake: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err = c.GetStake(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

// ListDelegations returns delegations for the given party.
func (n *DataNode) ListDelegations(req *dataapipb.ListDelegationsRequest) (response *dataapipb.ListDelegationsResponse, err error) {
	msg := "gRPC call failed (data-node): ListDelegations: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err = c.ListDelegations(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

type AccountFunds struct {
	Balance     *big.Int
	PartyId     string
	AccountType vega.AccountType
	AssetId     string
}

func (n *DataNode) ListAccounts(ctx context.Context, partyID string, accountType vega.AccountType, assetID *string) ([]AccountFunds, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	assetIDFilter := ""
	if assetID != nil {
		assetIDFilter = *assetID
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	requestCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()
	response, err := c.ListAccounts(requestCtx, &dataapipb.ListAccountsRequest{
		Filter: &dataapipb.AccountFilter{
			PartyIds:     []string{partyID},
			AccountTypes: []vega.AccountType{accountType},
			AssetId:      assetIDFilter,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var results []AccountFunds

	if response.Accounts == nil || len(response.Accounts.Edges) < 1 {
		return results, nil
	}

	for _, row := range response.Accounts.Edges {
		if row.Node == nil {
			continue
		}
		balance, _ := big.NewInt(0).SetString(row.Node.Balance, 10)

		results = append(results, AccountFunds{
			Balance:     balance,
			PartyId:     partyID,
			AssetId:     row.Node.Asset,
			AccountType: row.Node.Type,
		})
	}

	return results, nil
}
