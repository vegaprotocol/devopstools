package core

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// CoreClient stores state for a Vega Core node or Data Node.
type CoreClient struct {
	hosts       []string // format: host:port
	CallTimeout time.Duration
	Conn        *grpc.ClientConn
	mu          sync.RWMutex
	wg          sync.WaitGroup
	once        sync.Once
	logger      *zap.Logger
}

// NewCoreClient returns a new node.
func NewCoreClient(hosts []string, callTimeout time.Duration, logger *zap.Logger) *CoreClient {
	return &CoreClient{
		hosts:       hosts,
		CallTimeout: callTimeout,
		logger:      logger,
	}
}

// MustDialConnection tries to establish a connection to one of the nodes from a list of locations.
// It is idempotent, while it each call will block the caller until a connection is established.
func (n *CoreClient) MustDialConnection(ctx context.Context) {
	n.mustDialConnection(ctx, false)
}

func (n *CoreClient) MustDialConnectionIgnoreTime(ctx context.Context) {
	n.mustDialConnection(ctx, true)
}

func (n *CoreClient) mustDialConnection(ctx context.Context, ignoreTime bool) {
	n.once.Do(func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		n.wg.Add(len(n.hosts))

		for _, h := range n.hosts {
			go func(host string, ignoreTime bool) {
				defer n.wg.Done()
				if err := n.dialNode(ctx, host, ignoreTime); err == nil {
					cancel()
				} else {
					n.logger.Sugar().Debugf("failed to dial node(%s): %s", host, err.Error())
				}
			}(h, ignoreTime)
		}
		n.wg.Wait()
		n.mu.Lock()
		defer n.mu.Unlock()

		if n.Conn == nil {
			log.Fatalf("Failed to connect to DataNode")
		}
	})

	n.wg.Wait()
	n.once = sync.Once{}
}

func (n *CoreClient) dialNode(ctx context.Context, host string, ignoreTime bool) error {
	n.logger.Debug("dialing gRPC node", zap.String("node", host))
	conn, err := grpc.DialContext(
		ctx,
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		if err != context.Canceled {
			n.logger.Debug("Failed to dial node", zap.String("node", host), zap.Error(err))
		}
		return err
	}
	if conn.GetState() != connectivity.Ready {
		n.logger.Debug("Connection not ready", zap.String("node", host))
		return fmt.Errorf("connection not ready")
	}

	client := vegaapipb.NewCoreServiceClient(conn)
	res, err := client.Statistics(ctx, &vegaapipb.StatisticsRequest{})
	if err != nil {
		n.logger.Debug("Failed to get statistics", zap.String("node", host))
		return err
	}
	currentTime, err := time.Parse(time.RFC3339, res.Statistics.CurrentTime)
	if err != nil {
		return fmt.Errorf("failed to parse current time from statistics response %w", err)
	}
	vegaTime, err := time.Parse(time.RFC3339, res.Statistics.VegaTime)
	if err != nil {
		return fmt.Errorf("failed to parse vega time from statistics response %w", err)
	}

	if !ignoreTime && currentTime.Sub(vegaTime) > time.Minute*2 {
		n.logger.Debug(
			"node time is too far back",
			zap.String("node", host),
			zap.Time("vegaTime", vegaTime),
			zap.Time("currentTime", currentTime),
		)
		return fmt.Errorf("node time is too far back")
	}

	n.mu.Lock()
	n.Conn = conn
	n.mu.Unlock()
	return nil
}

func (n *CoreClient) Target() string {
	return n.Conn.Target()
}

func (n *CoreClient) WaitForStateChange(ctx context.Context, state connectivity.State) bool {
	return n.Conn.WaitForStateChange(ctx, state)
}
