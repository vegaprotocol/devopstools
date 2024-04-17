package ethereum

import (
	"context"
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20bridge"
	"github.com/vegaprotocol/devopstools/smartcontracts/multisigcontrol"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"

	"code.vegaprotocol.io/vega/protos/vega"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type ChainsClient struct {
	PrimaryChain *ChainClient
	EVMChain     *ChainClient
}

type ChainClient struct {
	logger *zap.Logger

	minterWallet     *Wallet
	collateralBridge *erc20bridge.ERC20Bridge
	multisigControl  *multisigcontrol.MultisigControl
	client           *ethclient.Client

	chainID string
}

func (c *ChainClient) EthClient() *ethclient.Client {
	return c.client
}

func (c *ChainClient) ID() string {
	return c.chainID
}

func NewChainClientForID(ctx context.Context, cfg config.Config, networkParams *types.NetworkParams, chainID string, logger *zap.Logger) (*ChainClient, error) {
	primaryEthConfig, err := networkParams.PrimaryEthereumConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get primary ethereum configuration from network paramters: %w", err)
	}

	evmChainConfig, err := networkParams.EVMChainConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get secondary ethereum configuration from network paramters: %w", err)
	}

	switch chainID {
	case primaryEthConfig.ChainId:
		primaryChainClient, err := NewPrimaryChainClient(ctx, cfg.Bridges.Primary, primaryEthConfig, logger.Named("primary-chain-client"))
		if err != nil {
			return nil, fmt.Errorf("could not initialize primary ethereum chain client: %w", err)
		}
		return primaryChainClient, nil
	case evmChainConfig.ChainId:
		evmChainClient, err := NewEVMChainClient(ctx, cfg.Bridges.EVM, evmChainConfig, logger.Named("evm-chain-client"))
		if err != nil {
			return nil, fmt.Errorf("could not initialize EVM chain client: %w", err)
		}
		return evmChainClient, nil
	default:
		return nil, fmt.Errorf("chain ID %q does not match any ethereum configuration in network parameter: %w", chainID, err)
	}
}

func NewChainClients(ctx context.Context, cfg config.Config, networkParams *types.NetworkParams, logger *zap.Logger) (*ChainsClient, error) {
	primaryEthConfig, err := networkParams.PrimaryEthereumConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get primary ethereum configuration from network paramters: %w", err)
	}

	evmChainConfig, err := networkParams.EVMChainConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get secondary ethereum configuration from network paramters: %w", err)
	}

	primaryChainClient, err := NewPrimaryChainClient(ctx, cfg.Bridges.Primary, primaryEthConfig, logger.Named("primary-chain-client"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize primary ethereum chain client: %w", err)
	}

	evmChainClient, err := NewEVMChainClient(ctx, cfg.Bridges.EVM, evmChainConfig, logger.Named("evm-chain-client"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize EVM chain client: %w", err)
	}

	return &ChainsClient{
		PrimaryChain: primaryChainClient,
		EVMChain:     evmChainClient,
	}, nil
}

func NewPrimaryChainClient(ctx context.Context, cfg config.PrimaryBridge, ethConfig *vega.EthereumConfig, logger *zap.Logger) (*ChainClient, error) {
	client, err := ethclient.DialContext(ctx, cfg.ClientURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ethereum client: %w", err)
	}

	collateralBridge, err := erc20bridge.NewERC20Bridge(client, ethConfig.CollateralBridgeContract.Address, erc20bridge.V2)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize collateral bridge client: %w", err)
	}

	multisigControl, err := multisigcontrol.NewMultisigControl(
		client,
		ethConfig.MultisigControlContract.Address,
		multisigcontrol.V2,
		cfg.Signers,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize multisig control client: %w", err)
	}

	minterWallet, err := tools.RetryReturn(6, 10*time.Second, func() (*Wallet, error) {
		w, err := NewWallet(ctx, client, cfg.Wallets.Minter.PrivateKey)
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
		multisigControl:  multisigControl,
		minterWallet:     minterWallet,

		chainID: ethConfig.ChainId,
	}, nil
}

func NewEVMChainClient(ctx context.Context, cfg config.EVMBridge, ethConfig *vega.EVMChainConfig, logger *zap.Logger) (*ChainClient, error) {
	client, err := ethclient.DialContext(ctx, cfg.ClientURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ethereum client: %w", err)
	}

	collateralBridge, err := erc20bridge.NewERC20Bridge(client, ethConfig.CollateralBridgeContract.Address, erc20bridge.V2)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize collateral bridge client: %w", err)
	}

	multisigControl, err := multisigcontrol.NewMultisigControl(
		client,
		ethConfig.MultisigControlContract.Address,
		multisigcontrol.V2,
		cfg.Signers,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize multisig contract: %w", err)
	}

	minterWallet, err := tools.RetryReturn(6, 10*time.Second, func() (*Wallet, error) {
		w, err := NewWallet(ctx, client, cfg.Wallets.Minter.PrivateKey)
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
		multisigControl:  multisigControl,
		chainID:          ethConfig.ChainId,
	}, nil
}
