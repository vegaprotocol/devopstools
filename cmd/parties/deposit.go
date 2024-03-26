package parties

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/smartcontracts/communitytoken"
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20bridge"

	vgfs "code.vegaprotocol.io/vega/libs/fs"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/polydawn/refmt/json"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DepositArgs struct {
	*Args
	EthereumPrivateKey      string
	ERC20TokenAddress       string
	CollateralBridgeAddress string
	DepositsFile            string
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
	depositCmd.PersistentFlags().StringVar(&depositArgs.CollateralBridgeAddress, "collateral-bridge-address", "", "The collateral bridge address")
	depositCmd.PersistentFlags().StringVar(&depositArgs.DepositsFile, "deposits-file", "deposits.json", "Path to the file containing the deposits as JSON: { \"<party>\": \"<amount>\" }")
}

func RunDepositToParties(args DepositArgs) error {
	ctx := context.Background()

	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return fmt.Errorf("failed to create vega network object: %w", err)
	}
	defer network.Disconnect()

	privateKey, err := crypto.HexToECDSA(args.EthereumPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decode private key from hex: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("failed casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// FIXME: On which criteria the selection should be done?
	ethClient := network.PrimaryEthClient

	nonce, err := ethClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get pending nonce: %w", err)
	}

	gasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get suggested gas price: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, nil)
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(800000) // in units
	auth.GasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(10))

	tokenAddress := common.HexToAddress(args.ERC20TokenAddress)

	communityTokenInstance, err := communitytoken.NewCommunitytoken(tokenAddress, ethClient)
	if err != nil {
		return err
	}
	symbol, _ := communityTokenInstance.Symbol(&bind.CallOpts{})
	name, _ := communityTokenInstance.Name(&bind.CallOpts{})
	fmt.Printf("TokenName: %s\nSymbol: %s\n", name, symbol)
	fmt.Printf("Approving")
	_ = communityTokenInstance
	collateralBridgeAddress := common.HexToAddress(args.CollateralBridgeAddress)
	tx, err := communityTokenInstance.Approve(auth, collateralBridgeAddress, big.NewInt(0).Mul(big.NewInt(90000000000), big.NewInt(1000000)))
	if err != nil {
		return err
	}
	fmt.Printf("Waiting for transaction to be minted: %s\n", tx.Hash().String())
	_, err = bind.WaitMined(ctx, ethClient, tx)
	if err != nil {
		return err
	}
	fmt.Println("Transaction sent")

	fmt.Printf("Depositing")
	collateralBridgeInstance, err := erc20bridge.NewERC20Bridge(ethClient, args.CollateralBridgeAddress, erc20bridge.ERC20BridgeV2)
	if err != nil {
		return fmt.Errorf("failed to create instance for collateral bridge: %w", err)
	}

	nonce, err = ethClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil
	}

	deposits, err := readDepositsFile(args.DepositsFile)
	if err != nil {
		return fmt.Errorf("could not read deposits file at %q: %w", args.DepositsFile, err)
	}

	initialNonce := nonce
	transactions := make([]*types.Transaction, 0, len(deposits))
	for party, amount := range deposits {
		tx, err := deposit(ctx, party, amount, privateKey, tokenAddress, ethClient, collateralBridgeInstance, initialNonce)
		if err != nil {
			return fmt.Errorf("could not desposit asset to party %q: %w", party, err)
		}
		transactions = append(transactions, tx)
		initialNonce = initialNonce + 1
	}

	for _, tx := range transactions {
		fmt.Printf("Waiting for transaction to be minted: %s\n", tx.Hash().String())
		if _, err := bind.WaitMined(ctx, ethClient, tx); err != nil {
			return err
		}

		fmt.Println("Transaction sent")
	}

	return nil
}

func deposit(ctx context.Context, partyId string, amount *big.Int, privateKey *ecdsa.PrivateKey, tokenAddress common.Address, client *ethclient.Client, bridgeInstance *erc20bridge.ERC20Bridge, nonce uint64) (*types.Transaction, error) {
	amount = amount.Mul(amount, big.NewInt(1000000))

	fmt.Printf("Deposit %s assets to party %s\n", amount.String(), partyId)

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, nil)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(800000) // in units
	auth.GasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(10))
	bytePubKey, err := hex.DecodeString(partyId)
	if err != nil {
		return nil, err
	}

	var party [32]byte
	copy(party[:], bytePubKey)
	tx, err := bridgeInstance.DepositAsset(auth, tokenAddress, amount, party)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Transaction sent: %v\n", tx)
	time.Sleep(1 * time.Second)

	return tx, nil
}

func readDepositsFile(path string) (map[string]*big.Int, error) {
	deposits := map[string]*big.Int{}

	content, err := vgfs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	if err := json.Unmarshal(content, &deposits); err != nil {
		return nil, fmt.Errorf("could not deserialize content: %w", err)
	}

	return deposits, nil
}
