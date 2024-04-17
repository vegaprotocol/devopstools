package core

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	e "github.com/vegaprotocol/devopstools/errors"

	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// maximumTimeDifference is the maximum
const maximumTimeDifference = 2 * time.Minute

// Client stores state for a Vega Core node or Data Node.
type Client struct {
	hosts       []string // format: host:port
	CallTimeout time.Duration
	Conn        *grpc.ClientConn
	mu          sync.RWMutex
	wg          sync.WaitGroup
	once        sync.Once
	logger      *zap.Logger
}

// NewClient returns a new node.
func NewClient(hosts []string, callTimeout time.Duration, logger *zap.Logger) *Client {
	return &Client{
		hosts:       hosts,
		CallTimeout: callTimeout,
		logger:      logger,
	}
}

// MustDialConnection tries to establish a connection to one of the nodes from a list of locations.
// It is idempotent, while it each call will block the caller until a connection is established.
func (n *Client) MustDialConnection(ctx context.Context) {
	n.mustDialConnection(ctx, false)
}

func (n *Client) MustDialConnectionIgnoreTime(ctx context.Context) {
	n.mustDialConnection(ctx, true)
}

func (n *Client) mustDialConnection(ctx context.Context, ignoreTime bool) {
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
					n.logger.Debug("Failed to dial node", zap.String("host", host), zap.Error(err))
				}
			}(h, ignoreTime)
		}
		n.wg.Wait()
		n.mu.Lock()
		defer n.mu.Unlock()

		if n.Conn == nil {
			log.Fatalf("Unable to connect to any configured node")
		}
	})

	n.wg.Wait()
	n.once = sync.Once{}
}

func (n *Client) dialNode(ctx context.Context, host string, ignoreTime bool) error {
	n.logger.Debug("Dialing gRPC node", zap.String("node", host))
	conn, err := grpc.DialContext(
		ctx,
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial node: %w", err)
	}
	if conn.GetState() != connectivity.Ready {
		return e.ErrConnectionNotReady
	}

	client := vegaapipb.NewCoreServiceClient(conn)
	res, err := client.Statistics(ctx, &vegaapipb.StatisticsRequest{})
	if err != nil {
		return fmt.Errorf("could not get node statistics: %w", err)
	}
	currentTime, err := time.Parse(time.RFC3339, res.Statistics.CurrentTime)
	if err != nil {
		return fmt.Errorf("failed to parse current time from statistics response %w", err)
	}
	vegaTime, err := time.Parse(time.RFC3339, res.Statistics.VegaTime)
	if err != nil {
		return fmt.Errorf("failed to parse vega time from statistics response %w", err)
	}

	if !ignoreTime && currentTime.Sub(vegaTime) > maximumTimeDifference {
		return fmt.Errorf("vega time is more than %s late compared to node time (vega time: %s, node time: %s)", maximumTimeDifference.String(), vegaTime, currentTime)
	}

	n.mu.Lock()
	n.Conn = conn
	n.mu.Unlock()
	return nil
}

func (n *Client) Target() string {
	return n.Conn.Target()
}

func (n *Client) WaitForStateChange(ctx context.Context, state connectivity.State) bool {
	return n.Conn.WaitForStateChange(ctx, state)
}
