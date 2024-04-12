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

// SubmitTransaction submits a signed v2 transaction.
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

// LastBlockData gets the latest blockchain data, height, hash and pow parameters.
func (n *Client) LastBlock(context.Context) (*vegaapipb.LastBlockHeightResponse, error) {
	msg := "gRPC call failed: LastBlockData: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()
	var response *vegaapipb.LastBlockHeightResponse
	response, err := c.LastBlockHeight(ctx, &vegaapipb.LastBlockHeightRequest{})
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return response, err
}

func (n *Client) Statistics() (*vegaapipb.StatisticsResponse, error) {
	msg := "gRPC call failed: Statistics: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()
	response, err := c.Statistics(ctx, &vegaapipb.StatisticsRequest{})
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return response, err
}

// PropagateChainEvent submits a signed v2 transaction.
func (n *Client) PropagateChainEvent(
	req *vegaapipb.PropagateChainEventRequest,
) (response *vegaapipb.PropagateChainEventResponse, err error) {
	msg := "gRPC call failed: PropagateChainEvent: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err = c.PropagateChainEvent(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *Client) CoreNetworkParameters(parameterKey string) ([]*vega.NetworkParameter, error) {
	if n == nil {
		return nil, fmt.Errorf("failed to get network parameters from the core api: %w", e.ErrNil)
	}

	if n.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(
			"failed to get network parameters from the core api: %w",
			e.ErrConnectionNotReady,
		)
	}

	c := vegaapipb.NewCoreStateServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	input := &vegaapipb.ListNetworkParametersRequest{}
	if parameterKey != "" {
		input.NetworkParameterKey = parameterKey
	}
	response, err := c.ListNetworkParameters(ctx, input)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get network parameters from the core api: %w",
			e.ErrorDetail(err),
		)
	}

	return response.NetworkParameters, nil
}
