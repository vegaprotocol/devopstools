package topup

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"code.vegaprotocol.io/vega/protos/vega"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
)

type TopUpWithTransferArgs struct {
	*TopUpArgs
	VegaNetworkName string
	TradersURL      string
}

var topUpWithTransferArgs TopUpWithTransferArgs

// topUpWithTransferCmd represents the traderbot command
var topUpWithTransferCmd = &cobra.Command{
	Use:   "with-transfer",
	Short: "TopUp parties on network with vega transfer",
	Long:  `TopUp parties on network with vega transfer`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunTopUpWithTransfer(topUpWithTransferArgs); err != nil {
			traderbotArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	topUpWithTransferArgs.TopUpArgs = &topUpArgs

	TopUpCmd.AddCommand(topUpWithTransferCmd)
	topUpWithTransferCmd.PersistentFlags().StringVar(&topUpWithTransferArgs.VegaNetworkName, "network", "", "Vega Network name")
	topUpWithTransferCmd.PersistentFlags().StringVar(&topUpWithTransferArgs.TradersURL, "traders-url", "", "Traders URL")
	if err := topUpWithTransferCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := topUpWithTransferCmd.MarkPersistentFlagRequired("traders-url"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

type Trader struct {
	Name       string `json:"name"`
	PublicKey  string `json:"pubKey"`
	Parameters struct {
		Base                      string  `json:"marketBase"`
		Quote                     string  `json:"marketQuote"`
		SettlementContractAddress string  `json:"marketSettlementEthereumContractAddress"`
		SettlementVegaAssetID     string  `json:"marketSettlementVegaAssetID"`
		WantedTokens              float64 `json:"wanted_tokens"`
		CurrentBalance            float64 `json:"balance"`
	} `json:"parameters"`
}

type TraderList map[string]Trader

const DefaultTopUpValue = 10000

type AssetTopUp struct {
	Symbol          string
	ContractAddress string
	VegaAssetId     string
	TotalAmount     *big.Int
	Parties         map[string]*big.Int
}

func RunTopUpWithTransfer(args TopUpWithTransferArgs) error {
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return fmt.Errorf("failed to create vega network object: %w", err)
	}
	defer network.Disconnect()
	networkAssets, err := network.DataNodeClient.GetAssets()
	if err != nil {
		return fmt.Errorf("failed to get assets from datanode: %w", err)
	}

	traders, err := tools.RetryReturn(10, 5*time.Second, func() (*TraderList, error) {
		return readTraders(args.TradersURL)
	})
	if err != nil {
		return fmt.Errorf("failed to read traders: %w", err)
	}

	topUpRegistry, err := determineTopUpAmount(networkAssets, *traders)
	if err != nil {
		return fmt.Errorf("failed to determine top up amounts for the assets: %w", err)
	}

	fmt.Printf("%v", topUpRegistry)

	return nil
}

func determineTopUpAmount(assets map[string]*vega.AssetDetails, traders TraderList) (map[string]AssetTopUp, error) {
	topUpRegistry := map[string]AssetTopUp{}

	for _, traderDetails := range traders {
		assetDetails, assetExists := assets[traderDetails.Parameters.SettlementVegaAssetID]
		if !assetExists {
			return nil, fmt.Errorf(
				"failed to find asset on network: bot is using the %s asset but it does not exist on the network",
				traderDetails.Parameters.SettlementVegaAssetID,
			)
		}

		if _, registryExists := topUpRegistry[traderDetails.Parameters.SettlementVegaAssetID]; !registryExists {
			topUpRegistry[traderDetails.Parameters.SettlementVegaAssetID] = AssetTopUp{
				Symbol:          assetDetails.Symbol,
				ContractAddress: traderDetails.Parameters.SettlementContractAddress,
				VegaAssetId:     traderDetails.Parameters.SettlementVegaAssetID,
				TotalAmount:     big.NewInt(0),
				Parties:         map[string]*big.Int{},
			}
		}

		currentEntry := topUpRegistry[traderDetails.Parameters.SettlementVegaAssetID]

		requiredTopUp := necessaryTopUp(
			traderDetails.Parameters.CurrentBalance,
			traderDetails.Parameters.WantedTokens,
			3.0,
		)

		if requiredTopUp == 0 {
			continue
		}

		topUpAmountWithZeros := big.NewInt(0).
			Mul(
				big.NewInt(int64(requiredTopUp)),
				big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(assetDetails.Decimals)), nil),
			)
		currentEntry.Parties[traderDetails.PublicKey] = topUpAmountWithZeros

		currentEntry.TotalAmount = big.NewInt(0).
			Add(
				topUpRegistry[traderDetails.Parameters.SettlementVegaAssetID].TotalAmount,
				topUpAmountWithZeros,
			)
	}

	return topUpRegistry, nil
}

func necessaryTopUp(currentBalance, wantedBalance, factor float64) float64 {
	if wantedBalance < 0.01 || wantedBalance > currentBalance {
		return currentBalance * factor
	}

	// top up not required
	return 0
}

func readTraders(tradersURL string) (*TraderList, error) {
	tr := &http.Transport{
		IdleConnTimeout:       30 * time.Second,
		DisableCompression:    true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Get(tradersURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send get request to %s: %w", tradersURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received invalid http code from the %s: got %d, expected 200", tradersURL, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := struct {
		Traders TraderList `json:"traders"`
	}{}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal traders response: %w", err)
	}

	return &response.Traders, nil
}
