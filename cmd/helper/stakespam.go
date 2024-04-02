package helper

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/vegaprotocol/devopstools/smartcontracts/erc20bridge"
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20token"
	"github.com/vegaprotocol/devopstools/smartcontracts/stakingbridge"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

type StakeSpamArgs struct {
	EthereumNodeAddress      string
	EthereumPrivateKey       string
	VegaAddressStakeReceiver string

	ERC20TokenAddress       string
	StakingBridgeAddress    string
	CollateralBridgeAddress string

	DoDeposit bool
	DoStake   bool

	RepeatAmount int
}

var stakeSpamArgs StakeSpamArgs

// selfDelegateCmd represents the selfDelegate command
var stakeSpamCmd = &cobra.Command{
	Use:   "spam-stake",
	Short: "Continously send stake orders",
	Long: `Let's say We have some tokens on the ethereum address assigned to the private key <PRIV-KEY>, this command can stake or deposit them on given bridges.
	The receiver is given vega address <VEGA-RECEIVER>.
	
	Example:
	./devopstools helper spam-stake \
	--ethereum-node "https://sepolia.infura.io/v3/XXXXXXXX" \
	--ethereum-private-key <PRIV-KEY> \
	--erc20-token-address "0x4d2f52bf29aae53f3bb0473e06c425d495a1ef76" \
	--vega-address-receiver "<VEGA-RECEIVER>" \
	--staking-bridge-address "0x7183c92cfDf82b22A704f4d34a49E64e8dF8580e" \
	--collateral-bridge-address "0x50BCd741E84EDebC155149d086a8A4BCBB878805"`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunStakeSpam(stakeSpamArgs); err != nil {
			log.Fatalf("Failed to execute command: %s", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	HelperCmd.AddCommand(stakeSpamCmd)

	stakeSpamCmd.PersistentFlags().StringVar(&stakeSpamArgs.EthereumNodeAddress, "ethereum-node", "", "The ethereum node address, e.g. https://sepolia.infura.io/v3/....")
	stakeSpamCmd.PersistentFlags().StringVar(&stakeSpamArgs.EthereumPrivateKey, "ethereum-private-key", "", "The ethereum private key, you want to send transactions from")
	stakeSpamCmd.PersistentFlags().StringVar(&stakeSpamArgs.ERC20TokenAddress, "erc20-token-address", "", "The ERC20 token address")
	stakeSpamCmd.PersistentFlags().StringVar(&stakeSpamArgs.VegaAddressStakeReceiver, "vega-address-receiver", "", "The vega address receiver for token stake/deposit")
	stakeSpamCmd.PersistentFlags().StringVar(&stakeSpamArgs.StakingBridgeAddress, "staking-bridge-address", "", "The staking bridge address")
	stakeSpamCmd.PersistentFlags().StringVar(&stakeSpamArgs.CollateralBridgeAddress, "collateral-bridge-address", "", "The collateral bridge address")
	stakeSpamCmd.PersistentFlags().BoolVar(&stakeSpamArgs.DoDeposit, "deposit", true, "If true, vega tokens are deposited to given receiver")
	stakeSpamCmd.PersistentFlags().BoolVar(&stakeSpamArgs.DoStake, "stake", false, "If true, vega tokens are staked to given receiver")
	stakeSpamCmd.PersistentFlags().IntVar(&stakeSpamArgs.RepeatAmount, "repeat-amount", 10, "Number of deposits(or stakes)")
}

func increaseAllowance(client *ethclient.Client, fromAddress common.Address, privateKey *ecdsa.PrivateKey, tokenAddress string, spenderAddress string) error {
	log.Printf("Increasing allowance for %s on %s token", spenderAddress, tokenAddress)
	if client == nil {
		return fmt.Errorf("client is nil")
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get pending nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get suggested gas price: %w", err)
	}

	tokenInstance, err := erc20token.NewERC20Token(client, tokenAddress, erc20token.ERC20TokenBase)
	if err != nil {
		return fmt.Errorf("failed to get instance of erc20 token: %w", err)
	}
	if tokenInstance == nil {
		return fmt.Errorf("erc20 token instance is invalid")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, nil)
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	transaction, err := tokenInstance.IncreaseAllowance(auth, common.HexToAddress(spenderAddress), big.NewInt(1000000000))
	if err != nil {
		return fmt.Errorf("failed to increase allowance: %w", err)
	}
	log.Printf("Waiting for transaction %s to be minted", transaction.Hash())
	bind.WaitMined(context.Background(), client, transaction)
	log.Printf("Transaction %s minted", transaction.Hash())

	return nil
}

func RunStakeSpam(args StakeSpamArgs) error {
	if args.DoDeposit && args.CollateralBridgeAddress == "" {
		return fmt.Errorf("when the --deposit flag is provided, collateral bridge address cannot be empty")
	}

	if args.DoStake && args.StakingBridgeAddress == "" {
		return fmt.Errorf("when the --deposit flag is provided, collateral bridge address cannot be empty")
	}

	if args.VegaAddressStakeReceiver == "" {
		return fmt.Errorf("vega address for receiver cannot be empty")
	}

	if args.EthereumPrivateKey == "" {
		return fmt.Errorf("ethereum private address cannot be empty")
	}

	if !args.DoDeposit && !args.DoStake {
		return fmt.Errorf("you have to provide at least one of the following flags: `--stake`, `--deposit`")
	}

	if args.RepeatAmount < 1 {
		return fmt.Errorf("repeat-amount must be greater than 0")
	}

	client, err := ethclient.Dial(args.EthereumNodeAddress)
	if err != nil {
		return fmt.Errorf("failed to create ethereum client: %w", err)
	}

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

	if args.DoStake {
		if err := increaseAllowance(client, fromAddress, privateKey, args.ERC20TokenAddress, args.StakingBridgeAddress); err != nil {
			return fmt.Errorf("failed to increase allowance for staking bridge: %w", err)
		}
	}

	if args.DoDeposit {
		if err := increaseAllowance(client, fromAddress, privateKey, args.ERC20TokenAddress, args.CollateralBridgeAddress); err != nil {
			return fmt.Errorf("failed to increase allowance for collateral bridge: %w", err)
		}
	}

	for i := 0; i < args.RepeatAmount; i++ {
		if args.DoDeposit {
			log.Println("Doing deposit")
			if err := executeSpam(client, fromAddress, privateKey, args.ERC20TokenAddress, args.VegaAddressStakeReceiver, args.CollateralBridgeAddress, ""); err != nil {
				log.Printf("Failed to deposit: %s", err.Error())
			}
		}

		if args.DoStake {
			log.Println("Doing stake")
			if err := executeSpam(client, fromAddress, privateKey, args.ERC20TokenAddress, args.VegaAddressStakeReceiver, "", args.StakingBridgeAddress); err != nil {
				log.Printf("Failed to stake: %s", err.Error())
			}
		}

		log.Println("... done")
	}

	return nil
}

func executeSpam(client *ethclient.Client,
	fromAddress common.Address,
	privateKey *ecdsa.PrivateKey,
	erc20TokenAddress string,
	vegaAddressTokenReceiver string,
	collateralBridgeAddress string,
	stakingBridgeAddress string,
) error {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get pending nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get suggested gas price: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, nil)
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	var transaction *types.Transaction
	if collateralBridgeAddress != "" {
		collateralBridgeInstance, err := erc20bridge.NewERC20Bridge(client, collateralBridgeAddress, erc20bridge.ERC20BridgeV2)
		if err != nil {
			return fmt.Errorf("failed to create instance for collateral bridge: %w", err)
		}

		vegaDepositReceiver := [32]byte{}
		depositReceiverSlice := []byte(vegaAddressTokenReceiver)
		for i := 0; i < 32; i++ {
			vegaDepositReceiver[i] = depositReceiverSlice[i]
		}
		transaction, err = collateralBridgeInstance.DepositAsset(auth, common.HexToAddress(erc20TokenAddress), big.NewInt(10), vegaDepositReceiver)
		if err != nil {
			return fmt.Errorf("failed to deposit: %w", err)
		}
	} else if stakingBridgeAddress != "" {
		stakingBridgeInstance, err := stakingbridge.NewStakingBridge(client, stakingBridgeAddress, stakingbridge.StakingBridgeV1)
		if err != nil {
			return fmt.Errorf("failed to create instance for staking bridge: %w", err)
		}

		transaction, err = stakingBridgeInstance.Stake(auth, big.NewInt(10), vegaAddressTokenReceiver)
		if err != nil {
			return fmt.Errorf("failed to stake: %w", err)
		}
	}
	if transaction == nil {
		return fmt.Errorf("empty transaction")
	}
	log.Printf("Waiting for transaction %s to be minted", transaction.Hash())
	bind.WaitMined(context.Background(), client, transaction)
	log.Printf("Transaction %s minted", transaction.Hash())

	return nil
}
