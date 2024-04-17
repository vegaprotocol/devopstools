package parties

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
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	vgfs "code.vegaprotocol.io/vega/libs/fs"

	"github.com/polydawn/refmt/json"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DepositArgs struct {
	*Args
	EthereumPrivateKey string
	ERC20TokenAddress  string
	DepositsFile       string
}

var depositArgs DepositArgs

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit tokens to given parties",
	Long:  "Deposit tokens to given parties",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunDepositToParties(depositArgs); err != nil {
			depositArgs.Logger.Error("An error occurred while depositing tokens to parties", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	depositArgs.Args = &args

	Cmd.AddCommand(depositCmd)
	depositCmd.PersistentFlags().StringVar(&depositArgs.EthereumPrivateKey, "ethereum-private-key", "", "The ethereum private key, you want to send transactions from")
	depositCmd.PersistentFlags().StringVar(&depositArgs.ERC20TokenAddress, "erc20-token-address", "", "The ERC20 token address")
	depositCmd.PersistentFlags().StringVar(&depositArgs.DepositsFile, "deposits-file", "deposits.json", "Path to the file containing the deposits as JSON: { \"<party>\": \"<amount>\" }")

	if err := depositCmd.MarkPersistentFlagRequired("ethereum-private-key"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := depositCmd.MarkPersistentFlagRequired("erc20-token-address"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := depositCmd.MarkPersistentFlagRequired("deposits-file"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunDepositToParties(args DepositArgs) error {
	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(args.NetworkFile)
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
	networkParams, err := datanodeClient.GetAllNetworkParameters()
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}
	logger.Debug("Network parameters retrieved")

	asset, erc20details, err := datanodeClient.ERC20AssetWithAddress(ctx, args.ERC20TokenAddress)
	if err != nil {
		return fmt.Errorf("failed to retrieve ERC20 asset with address %q: %w", args.ERC20TokenAddress, err)
	}

	chainClient, err := ethereum.NewChainClientForID(ctx, cfg, networkParams, erc20details.ChainId, logger)
	if err != nil {
		return fmt.Errorf("could not load chain client: %w", err)
	}

	deposits, err := readDepositsFile(args.DepositsFile, asset.Decimals)
	if err != nil {
		return fmt.Errorf("failed to parse deposits file at %q: %w", args.DepositsFile, err)
	}

	if err := chainClient.DepositERC20AssetFromAddress(ctx, args.EthereumPrivateKey, args.ERC20TokenAddress, deposits); err != nil {
		return fmt.Errorf("failed to deposit asset: %w", err)
	}

	logger.Info("Assets deposited successfully",
		zap.String("asset-name", asset.Name),
		zap.String("asset-symbol", asset.Symbol),
		zap.String("asset-contract-address", args.ERC20TokenAddress),
	)

	return nil
}

func readDepositsFile(path string, decimals uint64) (map[string]*types.Amount, error) {
	depositsAsSubUnit := map[string]*big.Int{}

	content, err := vgfs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	if err := json.Unmarshal(content, &depositsAsSubUnit); err != nil {
		return nil, fmt.Errorf("could not deserialize content: %w", err)
	}

	deposits := map[string]*types.Amount{}
	for partyID, amountAsSubUnit := range depositsAsSubUnit {
		deposits[partyID] = types.NewAmountFromSubUnit(amountAsSubUnit, decimals)
	}

	return deposits, nil
}
