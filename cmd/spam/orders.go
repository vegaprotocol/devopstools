package spam

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	"go.uber.org/zap"

	vegapb "code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"github.com/vegaprotocol/devopstools/wallet"
)

type OrdersArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
	MarketID        string
	Orders          uint64
	Threads         uint8
	MaxPrice        uint64
	MinPrice        uint64
	RateLimit       uint64
}

var ordersArgs OrdersArgs

var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Send a lot of orders to the market which stays in the book, but not trades",
	Example: `
# Start 4 threads, each thread sends 15 orders per seconds
devopstools spam orders \
    --network stagnet1 \
    --market-id 86948f946a64e14bb2e284f825cd46879d1cb581ce405cc62e4f74fcded190d3 \
    --threads 5 \
	--max-price 10000 \
    --thread-rate-limit 15 \
    --github-token secret-token	
`,
	Run: func(cmd *cobra.Command, args []string) {
		RunMarketSpam(ordersArgs)
	},
}

type ordersStats struct {
	sentOrders    uint64
	successOrders uint64
}

func (stat ordersStats) AsString() string {
	return fmt.Sprintf("SentOrders: %d, SuccessOrders: %d", stat.sentOrders, stat.successOrders)
}

type OrderSpammer struct {
	logger *zap.Logger

	dataNodeClient vegaapi.DataNodeClient
	dataNodeMutex  sync.RWMutex
	lastBlockData  *vegaapipb.LastBlockHeightResponse
	rootVegaWallet *wallet.VegaWallet
	spammerWallets []*wallet.VegaWallet

	lastBlockMonitorDone chan bool
	spammerDone          chan bool

	marketDetails *vegapb.Market

	threads   uint8
	marketId  string
	minPrice  uint64
	maxPrice  uint64
	rateLimit uint64

	stats      []ordersStats
	statsMutex sync.RWMutex
}

func NewOrderSpammer(threads uint8, marketId string, rateLimit, minPrice, maxPrice uint64, logger *zap.Logger, network *veganetwork.VegaNetwork) (*OrderSpammer, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if threads < 1 {
		return nil, fmt.Errorf("at least one thread must be running for spammer")
	}

	if network == nil {
		return nil, fmt.Errorf("network cannot be nil")
	}

	if minPrice > maxPrice {
		newMaxPrice := minPrice
		minPrice = maxPrice
		maxPrice = newMaxPrice
	}

	if maxPrice < 1 {
		return nil, fmt.Errorf("max price must be bigger than 0")
	}

	marketDetails, err := network.DataNodeClient.GetMarketById(marketId)
	if err != nil {
		return nil, fmt.Errorf("failed to get market details for spammer")
	}

	return &OrderSpammer{
		logger:    logger,
		threads:   threads,
		marketId:  marketId,
		minPrice:  minPrice,
		maxPrice:  maxPrice,
		rateLimit: rateLimit,

		spammerDone:          make(chan bool, threads),
		lastBlockMonitorDone: make(chan bool),

		marketDetails: marketDetails,

		rootVegaWallet: network.VegaTokenWhale,
		spammerWallets: make([]*wallet.VegaWallet, threads),
		dataNodeClient: network.DataNodeClient,
		stats:          make([]ordersStats, threads),
	}, nil
}

func (spammer *OrderSpammer) LastBlockMonitor() {
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	for {
		select {
		case <-spammer.lastBlockMonitorDone:
			spammer.logger.Info("DatanodeHeightMonitor stopped")
			return
		case <-ticker.C:
			lastBlockData, err := spammer.dataNodeClient.LastBlockData()
			if err != nil {
				spammer.logger.Error("failed to get data node last block data", zap.Error(err))
				// spammer.dataNodeClient.
				time.Sleep(time.Second)
				continue
			}
			spammer.dataNodeMutex.Lock()
			spammer.lastBlockData = lastBlockData
			spammer.dataNodeMutex.Unlock()
		}
	}
}

func getOrder(reference, pubKey, marketId string, minPrice, maxPrice uint64) *walletpb.SubmitTransactionRequest {
	price := minPrice + rand.Uint64()%(maxPrice-minPrice)

	return &walletpb.SubmitTransactionRequest{
		PubKey: pubKey,
		Command: &walletpb.SubmitTransactionRequest_OrderSubmission{
			OrderSubmission: &commandspb.OrderSubmission{
				MarketId:    marketId,
				Size:        1 + rand.Uint64()%5,
				Price:       fmt.Sprintf("%d", price),
				Side:        vegapb.Side_SIDE_BUY,
				TimeInForce: vegapb.Order_TIME_IN_FORCE_GTC,
				Type:        vegapb.Order_TYPE_LIMIT,
				Reference:   reference,
			},
		},
	}
}

func (spammer *OrderSpammer) Report() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-spammer.lastBlockMonitorDone:
			spammer.logger.Info("DatanodeHeightMonitor stopped")
			return
		case <-ticker.C:
			spammer.statsMutex.RLock()
			for thread, threadStats := range spammer.stats {
				spammer.logger.Info(
					fmt.Sprintf("Spam statistics: %s", threadStats.AsString()),
					zap.Int("thread", thread),
					zap.String("party", spammer.spammerWallets[thread].PublicKey),
				)
			}
			spammer.statsMutex.RUnlock()
		}
	}
}

func (spammer *OrderSpammer) Run() error {
	go spammer.LastBlockMonitor()

	for i := uint8(0); i < spammer.threads; i++ {
		vegaWallet, err := spammer.rootVegaWallet.DeriveKeyPair()
		if err != nil {
			return fmt.Errorf("failed to derive spammer wallet for thread %d", i)
		}

		go spammer.spammerThread(i, vegaWallet)
		spammer.spammerWallets[i] = vegaWallet
	}
	balanceManager, err := NewBalanceManager(spammer.dataNodeClient, spammer.rootVegaWallet, 60*time.Second)
	if err != nil {
		return fmt.Errorf("failed to create balance manager: %w", err)
	}

	go balanceManager.Run()

	for i := uint8(0); i < spammer.threads; i++ {
		assetPair := NewAssetPartyPair(
			spammer.marketDetails.TradableInstrument.Instrument.GetFuture().SettlementAsset,
			spammer.spammerWallets[i].PublicKey,
		)
		if err := balanceManager.AddAssetPartyPair(assetPair); err != nil {
			return fmt.Errorf("failed to add asset-party pair (%s) to balance manager: %w", assetPair.String(), err)
		}
	}

	go spammer.Report()

	return nil
}

func (spammer *OrderSpammer) spammerThread(idx uint8, wallet *wallet.VegaWallet) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(1000/(1+spammer.rateLimit)))
	defer ticker.Stop()

	spammer.logger.Info("Spammer thread starting", zap.Uint8("thread", idx), zap.Uint64("rate_limit", spammer.rateLimit), zap.Uint64("tick_every", 1000/(1+spammer.rateLimit)))
	for {
		select {
		case <-spammer.spammerDone:
			spammer.logger.Info("Spammer thread finished", zap.Uint8("thread", idx))
			return
		case <-ticker.C:
			if spammer.lastBlockData == nil {
				spammer.logger.Info("Spammer still getting required data from the network", zap.Uint8("thread", idx))
				time.Sleep(time.Second)
				continue
			}

			orderTx := getOrder(
				fmt.Sprintf("spammer-thread-%d", idx),
				wallet.PublicKey,
				spammer.marketId,
				spammer.minPrice,
				spammer.maxPrice,
			)

			spammer.dataNodeMutex.RLock()
			signedTx, err := wallet.SignTxWithPoW(orderTx, spammer.lastBlockData)
			spammer.dataNodeMutex.RUnlock()
			if err != nil {
				spammer.logger.Error("failed to sign transaction with pow", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}

			// wrap in vega Transaction Request
			submitReq := &vegaapipb.SubmitTransactionRequest{
				Tx:   signedTx,
				Type: vegaapipb.SubmitTransactionRequest_TYPE_SYNC,
			}

			spammer.statsMutex.Lock()
			spammer.stats[idx].sentOrders++
			spammer.statsMutex.Unlock()
			submitResponse, err := spammer.dataNodeClient.SubmitTransaction(submitReq)
			if err != nil {
				spammer.logger.Error("failed to send transaction", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}

			if !submitResponse.Success {
				spammer.logger.Error("order tranzaction failed", zap.String("log", submitResponse.Log), zap.String("data", submitResponse.Data))
				time.Sleep(time.Second)
				continue
			}

			spammer.statsMutex.Lock()
			spammer.stats[idx].successOrders++
			spammer.statsMutex.Unlock()
		}
	}
}

func (spammer *OrderSpammer) Stop() {
	spammer.lastBlockMonitorDone <- true

	for i := 0; i < int(spammer.threads); i++ {
		spammer.logger.Info("Stopping thread", zap.Int("thread", i))
		spammer.spammerDone <- true
	}
}

func RunMarketSpam(args OrdersArgs) error {
	rand.Seed(time.Now().UnixMicro())

	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return fmt.Errorf("failed to connect to the network: %w", err)
	}
	defer network.Disconnect()
	spammer, err := NewOrderSpammer(args.Threads, args.MarketID, args.RateLimit, args.MinPrice, args.MaxPrice, args.Logger, network)
	if err != nil {
		return fmt.Errorf("failed to create order spammer: %w", err)
	}

	spammer.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	spammer.Stop()
	time.Sleep(time.Second)

	return nil
}

func init() {
	ordersArgs.RootArgs = &rootCmd.Args

	ordersCmd.PersistentFlags().StringVar(&ordersArgs.VegaNetworkName, "network", "", "Vega Network name")
	ordersCmd.PersistentFlags().StringVar(&ordersArgs.MarketID, "market-id", "", "Market ID to spam the orders")
	ordersCmd.PersistentFlags().Uint8Var(&ordersArgs.Threads, "threads", 1, "Number of threads")
	ordersCmd.PersistentFlags().Uint64Var(&ordersArgs.MinPrice, "min-price", 1, "Minimum price")
	ordersCmd.PersistentFlags().Uint64Var(&ordersArgs.MaxPrice, "max-price", 10000, "Maximum price")
	ordersCmd.PersistentFlags().Uint64Var(&ordersArgs.RateLimit, "thread-rate-limit", 20, "The orders rate limit per second per thread")
	if err := ordersCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := ordersCmd.MarkPersistentFlagRequired("market-id"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
