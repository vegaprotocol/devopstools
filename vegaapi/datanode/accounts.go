package datanode

import (
	"context"
	"fmt"
	"math/big"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	e "github.com/vegaprotocol/devopstools/errors"
	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetPartyGeneralBalances(partyId string) (map[string]*big.Int, error) {
	res, err := n.ListAccounts(&dataapipb.ListAccountsRequest{
		Filter: &dataapipb.AccountFilter{
			PartyIds:     []string{partyId},
			AccountTypes: []vega.AccountType{vega.AccountType_ACCOUNT_TYPE_GENERAL},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get general balances for party %s, %w", partyId, err)
	}
	balanceByAssetId := map[string]*big.Int{}
	for _, edge := range res.Accounts.Edges {
		if edge.Node == nil {
			continue
		}
		assetId := edge.Node.Asset
		balance, ok := big.NewInt(0).SetString(edge.Node.Balance, 10)
		if !ok {
			return nil, fmt.Errorf("failed to get general balances for party %s, failed to parse balance %s", partyId, edge.Node.Balance)
		}
		if _, ok = balanceByAssetId[assetId]; ok {
			balanceByAssetId[assetId] = balanceByAssetId[assetId].Add(balanceByAssetId[assetId], balance)
		} else {
			balanceByAssetId[assetId] = balance
		}
	}

	return balanceByAssetId, nil
}

func (n *DataNode) ListAccounts(req *dataapipb.ListAccountsRequest) (response *dataapipb.ListAccountsResponse, err error) {
	msg := "gRPC call failed (data-node): ListAccounts: %w"
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

	response, err = c.ListAccounts(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}
