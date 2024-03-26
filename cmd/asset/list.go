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
	"github.com/vegaprotocol/devopstools/ethutils"

	"code.vegaprotocol.io/vega/libs/num"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type ListArgs struct {
	*Args
	AssetHexAddress   string
	VegaAssetID       string
	LifetimeLimit     string
	WithdrawThreshold string
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
	listCmd.PersistentFlags().StringVar(&listArgs.AssetHexAddress, "address", "", "Asset Hex Address")
	listCmd.PersistentFlags().StringVar(&listArgs.VegaAssetID, "vega-asset-id", "", "Vega asset ID")
	listCmd.PersistentFlags().StringVar(&listArgs.LifetimeLimit, "lifetime-limit", "", "Asset lifetime limit")
	listCmd.PersistentFlags().StringVar(&listArgs.WithdrawThreshold, "withdraw-threshold", "", "Asset withdraw threshold")

	if err := listCmd.MarkPersistentFlagRequired("vega-asset-id"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := listCmd.MarkPersistentFlagRequired("address"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := listCmd.MarkPersistentFlagRequired("lifetime-limit"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := listCmd.MarkPersistentFlagRequired("withdraw-threshold"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunListAssets(args ListArgs) error {
	ctx := context.Background()

	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return fmt.Errorf("failed to create vega network object: %w", err)
	}
	defer network.Disconnect()

	// parse Lifetime Limit
	bigLifetimeLimit, ok := new(big.Int).SetString(args.LifetimeLimit, 10)
	if !ok {
		return fmt.Errorf("failed to parse lifetime limit")
	}
	// parse Withdraw Threshold
	bigWithdrawThreshold, ok := new(big.Int).SetString(args.WithdrawThreshold, 10)
	if !ok {
		return fmt.Errorf("failed to parse withdraw threshold")
	}

	// TODO: how should we determine whether we should use primary or secondary bridge ?
	smartContract := network.PrimarySmartContracts
	ethClient := network.PrimaryEthClient
	multisigControl := network.PrimarySmartContracts.MultisigControl

	fmt.Printf("Listing %q asset to ERC20 Bridge for %q network... ", args.AssetHexAddress, args.VegaNetworkName)

	signersAddresses, err := multisigControl.GetSigners(ctx)
	if err != nil {
		return fmt.Errorf("failed to get signers addresses from multisig control: %w", err)
	}

	if len(signersAddresses) == 0 {
		return fmt.Errorf("no signers found for network %q", network.Network)
	}

	signersWallets := []*ethereum.EthWallet{}
	for _, node := range network.NodeSecrets {
		contained := slices.ContainsFunc(signersAddresses, func(address common.Address) bool {
			return address.Hex() == node.EthereumAddress
		})
		if !contained {
			continue
		}
		signerWallet, err := ethereum.NewWallet(context.Background(), ethClient, config.EthereumWallet{
			Address:    node.EthereumAddress,
			Mnemonic:   node.EthereumMnemonic,
			PrivateKey: node.EthereumPrivateKey,
		})
		if err != nil {
			return err
		}
		signersWallets = append(signersWallets, signerWallet)
	}

	assetIDB32, err := ethutils.VegaPubKeyToByte32(args.VegaAssetID)
	if err != nil {
		return err
	}

	assetSource := common.HexToAddress(args.AssetHexAddress)
	nonce := num.NewUint(ethutils.TimestampNonce()).BigInt()

	msg, err := listAssetMsg(smartContract.ERC20Bridge.Address, assetSource, assetIDB32, bigLifetimeLimit, bigWithdrawThreshold, nonce)
	if err != nil {
		return fmt.Errorf("could not generate message to sign: %w", err)
	}

	signatures, err := generateMultiSignature(signersWallets, msg)
	if err != nil {
		return fmt.Errorf("could not generate signatures: %w", err)
	}

	opts := network.NetworkMainWallet.GetTransactOpts(context.Background())
	tx, err := smartContract.ERC20Bridge.ListAsset(opts, assetSource, assetIDB32, bigLifetimeLimit, bigWithdrawThreshold, nonce, signatures)
	if err != nil {
		return err
	}

	if err = ethereum.WaitForTransaction(context.Background(), ethClient, tx, time.Minute); err != nil {
		return err
	}

	fmt.Printf("Success!\n")

	return nil
}

func generateMultiSignature(signers []*ethereum.EthWallet, msg []byte) ([]byte, error) {
	hash := crypto.Keccak256(msg)

	signatures := []byte{}
	for _, signerKey := range signers {
		signature, err := signerKey.Sign(hash)
		if err != nil {
			return nil, fmt.Errorf("failed to sign hash: %w", err)
		}

		signatures = append(signatures, signature...)
	}
	return signatures, nil
}

func listAssetMsg(bridgeAddr common.Address, tokenAddress common.Address, assetIDB32 [32]byte, lifetimeLimit *big.Int, withdrawThreshold *big.Int, nonce *big.Int) ([]byte, error) {
	typAddr, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, err
	}
	typBytes32, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, err
	}
	typString, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, err
	}
	typU256, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}

	args := abi.Arguments([]abi.Argument{
		{
			Name: "address",
			Type: typAddr,
		},
		{
			Name: "vega_asset_id",
			Type: typBytes32,
		},
		{
			Name: "lifetime_limit",
			Type: typU256,
		},
		{
			Name: "withdraw_treshold",
			Type: typU256,
		},
		{
			Name: "nonce",
			Type: typU256,
		},
		{
			Name: "func_name",
			Type: typString,
		},
	})

	buf, err := args.Pack([]interface{}{
		tokenAddress,
		assetIDB32,
		lifetimeLimit,
		withdrawThreshold,
		nonce,
		"list_asset",
	}...)
	if err != nil {
		return nil, fmt.Errorf("couldn't pack abi message: %w", err)
	}

	msg, err := packBufAndSubmitter(buf, bridgeAddr)
	if err != nil {
		return nil, fmt.Errorf("couldn't pack abi message: %w", err)
	}

	return msg, nil
}

func packBufAndSubmitter(buf []byte, submitter common.Address) ([]byte, error) {
	typBytes, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return nil, err
	}
	typAddr, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, err
	}

	args2 := abi.Arguments([]abi.Argument{
		{
			Name: "bytes",
			Type: typBytes,
		},
		{
			Name: "address",
			Type: typAddr,
		},
	})

	return args2.Pack(buf, submitter)
}
