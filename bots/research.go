package bots

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/wallet"
)

const (
	// marketMakerWalletIndex defines index of the market marker wallet. This is
	// hardcoded in the vega-market-sim.
	marketMakerWalletIndex = 3
)

type ResearchBot struct {
	TradingBot
	WalletData struct {
		Index          int64  `json:"index"`
		PublicKey      string `json:"publicKey"`
		RecoveryPhrase string `json:"recoveryPhrase"`
	} `json:"wallet"`

	wallet *wallet.VegaWallet
}

func (b *ResearchBot) IsMarketMaker() bool {
	return b.WalletData.Index != marketMakerWalletIndex
}

func (b *ResearchBot) GetWallet() (*wallet.VegaWallet, error) {
	if b.wallet == nil {
		w, err := wallet.GetVegaWalletSingleton(b.WalletData.RecoveryPhrase, uint32(b.WalletData.Index))
		if err != nil {
			return nil, fmt.Errorf("could not retrieve wallet: %w", err)
		}
		b.wallet = w
	}
	return b.wallet, nil
}
