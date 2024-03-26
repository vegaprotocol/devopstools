package ethereum

import (
	"context"
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20bridge"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"

	"code.vegaprotocol.io/vega/protos/vega"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type ChainsClient struct {
	PrimaryChain   *ChainClient
	SecondaryChain *ChainClient
}

type ChainClient struct {
	logger *zap.Logger

	minterWallet     *EthWallet
	collateralBridge *erc20bridge.ERC20Bridge
	client           *ethclient.Client

	chainID string
}

func (c *ChainClient) ID() string {
	return c.chainID
}

func NewChainClients(ctx context.Context, cfg config.Config, networkParams *types.NetworkParams, logger *zap.Logger) (*ChainsClient, error) {
	primaryEthConfig, err := networkParams.PrimaryEthereumConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get primary ethereum configuration from network paramters: %w", err)
	}

	secondaryEthConfig, err := networkParams.EVMChainConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get secondary ethereum configuration from network paramters: %w", err)
	}

	primaryChainClient, err := NewPrimaryChainClient(ctx, cfg.Bridges.Primary, primaryEthConfig, logger.Named("primary-chain-client"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize primary ethereum chain client: %w", err)
	}

	secondaryChainClient, err := NewSecondaryChainClient(ctx, cfg.Bridges.Secondary, secondaryEthConfig, logger.Named("primary-chain-client"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize primary ethereum chain client: %w", err)
	}

	return &ChainsClient{
		PrimaryChain:   primaryChainClient,
		SecondaryChain: secondaryChainClient,
	}, nil
}

func NewPrimaryChainClient(ctx context.Context, cfg config.PrimaryBridge, ethConfig *vega.EthereumConfig, logger *zap.Logger) (*ChainClient, error) {
	client, err := ethclient.DialContext(ctx, cfg.ClientURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ethereum client: %w", err)
	}

	collateralBridge, err := erc20bridge.NewERC20Bridge(client, ethConfig.CollateralBridgeContract.Address, erc20bridge.ERC20BridgeV2)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize collateral bridge client: %w", err)
	}

	minterWallet, err := tools.RetryReturn(6, 10*time.Second, func() (*EthWallet, error) {
		w, err := NewWallet(ctx, client, cfg.Wallets.Minter)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Ethereum wallet: %w", err)
		}
		return w, nil
	})
	if err != nil {
		return nil, err
	}

	return &ChainClient{
		logger: logger,

		client:           client,
		collateralBridge: collateralBridge,
		minterWallet:     minterWallet,

		chainID: ethConfig.ChainId,
	}, nil
}

func NewSecondaryChainClient(ctx context.Context, cfg config.SecondaryBridge, ethConfig *vega.EVMChainConfig, logger *zap.Logger) (*ChainClient, error) {
	client, err := ethclient.DialContext(ctx, cfg.ClientURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ethereum client: %w", err)
	}

	collateralBridge, err := erc20bridge.NewERC20Bridge(client, ethConfig.CollateralBridgeContract.Address, erc20bridge.ERC20BridgeV2)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize collateral bridge client: %w", err)
	}

	minterWallet, err := tools.RetryReturn(6, 10*time.Second, func() (*EthWallet, error) {
		w, err := NewWallet(ctx, client, cfg.Wallets.Minter)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Ethereum wallet: %w", err)
		}
		return w, nil
	})
	if err != nil {
		return nil, err
	}

	return &ChainClient{
		logger: logger,

		client:           client,
		collateralBridge: collateralBridge,
		minterWallet:     minterWallet,

		chainID: ethConfig.ChainId,
	}, nil
}
