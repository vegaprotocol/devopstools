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

func (n *DataNode) GetPartyTotalStake(ctx context.Context, partyId string) (*big.Int, error) {
	if n.Client.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Client.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.Client.CallTimeout)
	defer cancel()

	response, err := c.GetStake(reqCtx, &dataapipb.GetStakeRequest{
		PartyId: partyId,
	})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	totalStake := new(big.Int)
	totalStake, ok := totalStake.SetString(response.CurrentStakeAvailable, 0)
	if !ok {
		return nil, fmt.Errorf("failed to convert %s to big.Int", response.CurrentStakeAvailable)
	}

	return totalStake, nil
}

// ListDelegations returns delegations for the given party.
func (n *DataNode) ListDelegations(req *dataapipb.ListDelegationsRequest) (*dataapipb.ListDelegationsResponse, error) {
	if n.Client.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Client.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.Client.CallTimeout)
	defer cancel()

	response, err := c.ListDelegations(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	return response, nil
}

type AccountFunds struct {
	Balance     *big.Int
	PartyId     string
	AccountType vega.AccountType
	AssetId     string
}

func (n *DataNode) ListAccounts(ctx context.Context, partyID string, accountType vega.AccountType, assetID *string) ([]AccountFunds, error) {
	if n.Client.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	assetIDFilter := ""
	if assetID != nil {
		assetIDFilter = *assetID
	}

	c := dataapipb.NewTradingDataServiceClient(n.Client.Conn)
	requestCtx, cancel := context.WithTimeout(ctx, n.Client.CallTimeout)
	defer cancel()
	response, err := c.ListAccounts(requestCtx, &dataapipb.ListAccountsRequest{
		Filter: &dataapipb.AccountFilter{
			PartyIds:     []string{partyID},
			AccountTypes: []vega.AccountType{accountType},
			AssetId:      assetIDFilter,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
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
