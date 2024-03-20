package benchmark

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type EthRPCArgs struct {
	*BenchmarkArgs
	URLs   []string
	Repeat uint16
}

var ethRPCArgs EthRPCArgs

// ethRPCCmd represents the ethRPC command
var ethRPCCmd = &cobra.Command{
	Use:   "eth-endpoint",
	Short: "Benchmark Ethereum RPC endpoints",
	Long:  `Benchmark Ethereum RPC endpoints`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunEthRPC(ethRPCArgs); err != nil {
			ethRPCArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	ethRPCArgs.BenchmarkArgs = &benchmarkArgs

	BenchmarkCmd.AddCommand(ethRPCCmd)

	ethRPCCmd.PersistentFlags().StringSliceVar(&ethRPCArgs.URLs, "rpc", nil, "Comma separated list of Ethereum RPC endpoints")
	if err := ethRPCCmd.MarkPersistentFlagRequired("rpc"); err != nil {
		log.Fatalf("%v\n", err)
	}
	ethRPCCmd.PersistentFlags().Uint16Var(&ethRPCArgs.Repeat, "repeat", 1, "Repeat multiple times.")
}

func RunEthRPC(args EthRPCArgs) error {
	var err error
	endpoints := make([]RPCEndpoint, len(args.URLs))

	for i, url := range args.URLs {
		endpoints[i].URL = url
		endpoints[i].Client, err = ethclient.Dial(url)
		if err != nil {
			return fmt.Errorf("Failed to connect to %s: %w", url, err)
		}
	}

	for i := uint16(1); i <= args.Repeat; i++ {
		GetLatestBlockBenchmark(endpoints)
		SubscribeToEvents(endpoints)
	}

	return nil
}

type RPCEndpoint struct {
	URL    string
	Client *ethclient.Client
}

type LatestBlockResult struct {
	Endpoint RPCEndpoint
	Block    string
}

func GetLatestBlockBenchmark(endpoints []RPCEndpoint) {
	var wg sync.WaitGroup
	results := make(chan LatestBlockResult, len(endpoints))

	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(endpoint RPCEndpoint, wg *sync.WaitGroup, results chan LatestBlockResult) {
			defer wg.Done()
			header, err := endpoint.Client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Fatalf("Failed to get header for url %s, %s", endpoint.URL, err)
			}
			block := header.Number.String()
			results <- LatestBlockResult{
				Endpoint: endpoint,
				Block:    block,
			}
		}(endpoint, &wg, results)
	}
	wg.Wait()
	close(results)

	fmt.Print("Results:\n")
	for res := range results {
		fmt.Printf("- %s: %s\n", res.Endpoint.URL, res.Block)
	}
}

func SubscribeToEvents(endpoints []RPCEndpoint) {
	var wg sync.WaitGroup
	usdcContractAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	currentBlock, err := endpoints[0].Client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get block height for %s: %s", endpoints[0].URL, err)
	}
	usdcQuery := ethereum.FilterQuery{
		Addresses: []common.Address{usdcContractAddress},
		FromBlock: big.NewInt(int64(currentBlock - 10)),
	}

	fmt.Printf("Get events from block %d for USDC: %s\n", currentBlock-10, usdcContractAddress)

	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(endpoint RPCEndpoint) {
			defer wg.Done()
			logs, err := endpoint.Client.FilterLogs(context.Background(), usdcQuery)
			if err != nil {
				log.Fatalf("Failed on %s: %s", endpoint.URL, err)
			}
			evtNoPerBlock := make(map[uint64]int)
			for _, vLog := range logs {
				evtNoPerBlock[vLog.BlockNumber] += 1
				// fmt.Printf("%s: %s: %s\n", time.Now().Format(time.RFC850), endpoint.URL, vLog.TxHash)
			}

			sortedBlocks := make([]uint64, 0, len(evtNoPerBlock))
			for block := range evtNoPerBlock {
				sortedBlocks = append(sortedBlocks, block)
			}
			sort.Slice(sortedBlocks, func(i, j int) bool { return sortedBlocks[i] < sortedBlocks[j] })

			for _, block := range sortedBlocks {
				fmt.Printf("%s [%d]: %d\n", endpoint.URL, block, evtNoPerBlock[block])
			}
		}(endpoint)
	}
	wg.Wait()
}
