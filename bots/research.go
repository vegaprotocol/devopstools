package bots

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/vega"

	"code.vegaprotocol.io/vega/wallet/wallet"
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

	wallet wallet.Wallet
}

func (b *ResearchBot) IsMarketMaker() bool {
	return b.WalletData.Index != marketMakerWalletIndex
}

func (b *ResearchBot) GetWallet() (wallet.Wallet, error) {
	if b.wallet == nil {
		w, err := vega.LoadWalletAsSingleton(b.WalletData.RecoveryPhrase)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve singleton wallet: %w", err)
		}

		if err := vega.GenerateKeysUpToIndex(w, uint32(b.WalletData.Index)); err != nil {
			return nil, fmt.Errorf("could not generate keys: %w", err)
		}

		b.wallet = w
	}
	return b.wallet, nil
}
