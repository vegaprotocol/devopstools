package api

import (
	"context"
	"fmt"
	"math/big"

	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"code.vegaprotocol.io/vega/protos/vega"
)

type PartyClient struct {
	client *datanode.DataNode

	partyID string
}

func (a *PartyClient) PartyID() string {
	return a.partyID
}

func (a *PartyClient) GeneralAccountBalanceForAsset(ctx context.Context, assetID string) (*big.Int, error) {
	whaleAccounts, err := a.client.ListAccounts(ctx, a.partyID, vega.AccountType_ACCOUNT_TYPE_GENERAL, &assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve general accounts: %w", err)
	}

	whaleFundsAsSubUnits := big.NewInt(0)
	if len(whaleAccounts) > 0 {
		whaleFundsAsSubUnits = whaleAccounts[0].Balance
	}

	return whaleFundsAsSubUnits, nil
}

func NewPartyClient(client *datanode.DataNode, partyID string) *PartyClient {
	return &PartyClient{
		client: client,

		partyID: partyID,
	}
}
