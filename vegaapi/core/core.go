package core

import (
	"context"
	"fmt"

	e "github.com/vegaprotocol/devopstools/errors"

	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"

	"google.golang.org/grpc/connectivity"
)

func (n *Client) SendTransaction(ctx context.Context, tx *commandspb.Transaction, reqType vegaapipb.SubmitTransactionRequest_Type) (*vegaapipb.SubmitTransactionResponse, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	reqCtx, cancelReq := context.WithTimeout(ctx, n.CallTimeout)
	defer cancelReq()

	response, err := c.SubmitTransaction(reqCtx, &vegaapipb.SubmitTransactionRequest{
		Tx:   tx,
		Type: reqType,
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return response, nil
}

func (n *Client) LastBlock(ctx context.Context) (*vegaapipb.LastBlockHeightResponse, error) {

	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	response, err := c.LastBlockHeight(reqCtx, &vegaapipb.LastBlockHeightRequest{})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return response, nil
}

func (n *Client) Statistics(ctx context.Context) (*vegaapipb.StatisticsResponse, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	response, err := c.Statistics(reqCtx, &vegaapipb.StatisticsRequest{})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return response, nil
}

func (n *Client) CoreNetworkParameters(ctx context.Context, parameterKey string) ([]*vega.NetworkParameter, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := vegaapipb.NewCoreStateServiceClient(n.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	input := &vegaapipb.ListNetworkParametersRequest{}
	if parameterKey != "" {
		input.NetworkParameterKey = parameterKey
	}
	response, err := c.ListNetworkParameters(reqCtx, input)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	return response.NetworkParameters, nil
}
