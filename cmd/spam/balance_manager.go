package spam

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"

	vegapb "code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
)

const (
	topUpOvercomitPercentage = 150  // %
	wantedBalance            = 1000 // without decimals
)

type AssetPartyPair struct {
	assetId string
	partyId string

	wanted  *big.Int
	balance *big.Int
}

func (pair AssetPartyPair) String() string {
	return fmt.Sprintf(`{"AssetID": "%s", "PartyId": "%s"`, pair.assetId, pair.partyId)
}

func NewAssetPartyPair(assetId, partyId string) AssetPartyPair {
	return AssetPartyPair{
		assetId: assetId,
		partyId: partyId,
		wanted:  big.NewInt(0),
		balance: big.NewInt(0),
	}
}

type BalanceManager struct {
	dataNodeClient vegaapi.DataNodeClient
	whaleWallet    *wallet.VegaWallet

	assetsDecimals map[string]uint64

	tickPeriod      time.Duration
	assetPartyPairs []AssetPartyPair
}

func NewBalanceManager(dataNodeClient vegaapi.DataNodeClient, whaleWallet *wallet.VegaWallet, tickPeriod time.Duration) (*BalanceManager, error) {
	if whaleWallet == nil {
		return nil, fmt.Errorf("the whale wallet must be valid whale wallet")
	}

	if dataNodeClient == nil {
		return nil, fmt.Errorf("the data-node client cannot be nil")
	}

	if tickPeriod < 100*time.Millisecond {
		return nil, fmt.Errorf("tick period needs to be at least 100ms")
	}

	assets, err := describeAssets(dataNodeClient)
	if err != nil {
		return nil, fmt.Errorf("failed to describe assets: %w", err)
	}

	return &BalanceManager{
		dataNodeClient: dataNodeClient,
		whaleWallet:    whaleWallet,
		tickPeriod:     tickPeriod,
		assetsDecimals: assets,
	}, nil
}

func (bm *BalanceManager) AddAssetPartyPair(pair AssetPartyPair) error {
	assetDecimals, assetExists := bm.assetsDecimals[pair.assetId]
	if !assetExists {
		return fmt.Errorf("balance manager does not know decimals for asset %s", pair.assetId)
	}

	decimalsExp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(assetDecimals)), nil)
	pair.wanted = big.NewInt(0).Mul(big.NewInt(wantedBalance), decimalsExp)
	bm.assetPartyPairs = append(bm.assetPartyPairs, pair)

	return nil
}

func (bm *BalanceManager) TopUpPair(lastBlockData *vegaapipb.LastBlockHeightResponse, pair AssetPartyPair, amount *big.Int) error {
	if bm.whaleWallet == nil {
		return fmt.Errorf("the whale wallet must be valid whale wallet")
	}

	if lastBlockData == nil {
		return fmt.Errorf("data node details are nil")
	}

	log.Printf("adding %s token %s to party %s", amount.String(), pair.assetId, pair.partyId)

	transferTx := &walletpb.SubmitTransactionRequest{
		PubKey: bm.whaleWallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_Transfer{
			Transfer: &commandspb.Transfer{
				FromAccountType: vegapb.AccountType_ACCOUNT_TYPE_GENERAL,
				To:              pair.partyId,
				ToAccountType:   vegapb.AccountType_ACCOUNT_TYPE_GENERAL,
				Asset:           pair.assetId,
				Amount:          amount.String(),
				Reference:       "spammer-top-up",
				Kind: &commandspb.Transfer_OneOff{
					OneOff: &commandspb.OneOffTransfer{
						DeliverOn: time.Now().Add(time.Second * 2).Unix(),
					},
				},
			},
		},
	}

	signedTx, err := bm.whaleWallet.SignTxWithPoW(transferTx, lastBlockData)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// wrap in vega Transaction Request
	submitReq := &vegaapipb.SubmitTransactionRequest{
		Tx:   signedTx,
		Type: vegaapipb.SubmitTransactionRequest_TYPE_SYNC,
	}

	submitResponse, err := bm.dataNodeClient.SubmitTransaction(submitReq)
	if err != nil {
		return fmt.Errorf("failed to send the signed transaction: %w", err)
	}

	if !submitResponse.Success {
		return fmt.Errorf("sent transaction failed: %s", submitResponse.Data)
	}

	log.Printf("added %s token %s to party %s", amount.String(), pair.assetId, pair.partyId)

	return nil
}

func describeAssets(dataNodeClient vegaapi.DataNodeClient) (map[string]uint64, error) {
	assets, err := dataNodeClient.GetAssets()
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}

	results := map[string]uint64{}
	for id, asset := range assets {
		results[id] = asset.Decimals
	}

	return results, nil
}

func (bm *BalanceManager) requiredTopup(pair *AssetPartyPair) (*big.Int, error) {
	funds, err := bm.dataNodeClient.GetFunds(pair.partyId, vegapb.AccountType_ACCOUNT_TYPE_GENERAL, &pair.assetId)
	if err != nil {
		return nil, fmt.Errorf("failed to get funds for %s: %w", pair, err)
	}

	// given party has no funds at all
	if len(funds) < 1 {
		return pair.wanted, nil
	}
	if len(funds) > 1 {
		log.Printf("WARN: returned more than one funds for general account for party: %s", pair.String())
	}

	pair.balance = funds[0].Balance
	return big.NewInt(0).Sub(pair.wanted, funds[0].Balance), nil
}

func (bm *BalanceManager) Run() error {
	ticker := time.NewTicker(bm.tickPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lastBlockData, err := bm.dataNodeClient.LastBlockData()
			if err != nil {
				log.Printf("failed to get latest block info in balance manager: %s", err.Error())
				break
			}

			for _, pair := range bm.assetPartyPairs {
				requiredTopup, err := bm.requiredTopup(&pair)
				if err != nil {
					log.Printf("failed to check if topup required: %s", err.Error())
					continue
				}

				if requiredTopup.Cmp(big.NewInt(0)) < 1 {
					log.Printf("top up for party %s not required. available funds: %s", pair.partyId, pair.balance.String())
					continue
				}

				// top up = wanted * topUpOvercomitPercentage/100
				topUpOvercommited := big.NewInt(0).Mul(pair.wanted, big.NewInt(topUpOvercomitPercentage))
				topUpOvercommited = big.NewInt(0).Div(topUpOvercommited, big.NewInt(100))

				if err := bm.TopUpPair(lastBlockData, pair, topUpOvercommited); err != nil {
					log.Printf("failed to toup party \"%s\" with asset \"%s\": %s", pair.partyId, pair.assetId, err.Error())
				}
			}
		}
	}
}
