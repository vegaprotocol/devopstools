package asset

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type ListArgs struct {
	*Args
	ERC20TokenAddress string
	VegaAssetID       string
	LifetimeLimit     string
	WithdrawThreshold string
	ChainID           string
}

var listArgs ListArgs

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List an asset on the bridge",
	Long:  "List an asset on the bridge",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunListAssets(listArgs); err != nil {
			listArgs.Logger.Error("An error occurred while depositing tokens to parties", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	listArgs.Args = &args

	Cmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringVar(&listArgs.ERC20TokenAddress, "erc20-token-address", "", "Asset Hex Address")
	listCmd.PersistentFlags().StringVar(&listArgs.VegaAssetID, "vega-asset-id", "", "Vega asset ID")
	listCmd.PersistentFlags().StringVar(&listArgs.LifetimeLimit, "lifetime-limit", "", "Asset lifetime limit")
	listCmd.PersistentFlags().StringVar(&listArgs.WithdrawThreshold, "withdraw-threshold", "", "Asset withdraw threshold")
	listCmd.PersistentFlags().StringVar(&listArgs.ChainID, "chain-id", "", "Chain on which the contract is located")

	if err := listCmd.MarkPersistentFlagRequired("vega-asset-id"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := listCmd.MarkPersistentFlagRequired("erc20-token-address"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := listCmd.MarkPersistentFlagRequired("lifetime-limit"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := listCmd.MarkPersistentFlagRequired("withdraw-threshold"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := listCmd.MarkPersistentFlagRequired("chain-id"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunListAssets(args ListArgs) error {
	lifetimeLimit, ok := new(big.Int).SetString(args.LifetimeLimit, 10)
	if !ok {
		return fmt.Errorf("failed to parse lifetime limit")
	}
	withdrawThreshold, ok := new(big.Int).SetString(args.WithdrawThreshold, 10)
	if !ok {
		return fmt.Errorf("failed to parse withdraw threshold")
	}

	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(ctx, args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Debug("Network file loaded", zap.String("name", cfg.Name.String()))

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Debug("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Debug("Looking for healthy gRPC endpoints...")
	healthyEndpoints := tools.FilterHealthyGRPCEndpoints(endpoints)
	if len(healthyEndpoints) == 0 {
		return fmt.Errorf("no healthy gRPC endpoint found on configured datanodes")
	}
	logger.Debug("Healthy gRPC endpoints found", zap.Strings("endpoints", healthyEndpoints))

	datanodeClient := datanode.New(healthyEndpoints, 3*time.Second, args.Logger.Named("datanode"))

	logger.Debug("Connecting to a datanode's gRPC endpoint...")
	dialCtx, cancelDialing := context.WithTimeout(ctx, 2*time.Second)
	defer cancelDialing()
	datanodeClient.MustDialConnection(dialCtx)
	logger.Debug("Connected to a datanode's gRPC node", zap.String("node", datanodeClient.Target()))

	logger.Debug("Retrieving network parameters...")
	networkParams, err := datanodeClient.ListNetworkParameters(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}
	logger.Debug("Network parameters retrieved")

	chainClient, err := ethereum.NewChainClientForID(ctx, cfg, networkParams, args.ChainID, logger)
	if err != nil {
		return fmt.Errorf("could not load chain client: %w", err)
	}

	signersAddresses, err := chainClient.Signers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get signers addresses from multisig control: %w", err)
	}

	signersWallets, err := loadSignersWallets(ctx, cfg, signersAddresses, chainClient.EthClient())
	if err != nil {
		return err
	}

	if err := chainClient.ListAsset(ctx, signersWallets, args.VegaAssetID, args.ERC20TokenAddress, lifetimeLimit, withdrawThreshold); err != nil {
		return fmt.Errorf("failed to list asset: %w", err)
	}

	logger.Info("Asset listed successfully",
		zap.String("asset-id", args.VegaAssetID),
		zap.String("asset-contract-address", args.ERC20TokenAddress),
	)

	return nil
}

func loadSignersWallets(ctx context.Context, cfg config.Config, signersAddresses []common.Address, ethClient *ethclient.Client) ([]*ethereum.Wallet, error) {
	var signersWallets []*ethereum.Wallet
	for _, node := range cfg.Nodes {
		if !slices.ContainsFunc(signersAddresses, func(address common.Address) bool {
			return address.Hex() == node.Secrets.EthereumAddress
		}) {
			continue
		}

		signerWallet, err := ethereum.NewWallet(ctx, ethClient, node.Secrets.EthereumPrivateKey)
		if err != nil {
			return nil, err
		}
		signersWallets = append(signersWallets, signerWallet)
	}
	return signersWallets, nil
}
