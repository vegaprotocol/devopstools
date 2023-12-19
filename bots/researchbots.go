package bots

import (
	"fmt"
	"log"

	"github.com/vegaprotocol/devopstools/wallet"
)

const LeaderWalletIndex = 0

type ResearchBot struct {
	BotTrader
	WalletData struct {
		Index          int64   `json:"index"`
		PublicKey      string  `json:"publicKey"`
		RecoveryPhrase *string `json:"recoveryPhrase"`
	} `json:"wallet"`

	wallet *wallet.VegaWallet
}

type ResearchBots map[string]ResearchBot

func GetResearchBots(
	network string,
	botsAPIToken string,
) (ResearchBots, error) {
	botsURL := fmt.Sprintf("https://%s.bots.vega.rocks/traders", network)
	log.Printf("Getting research bot traders from: %s", botsURL)
	var payload struct {
		Traders map[string]ResearchBot `json:"traders"`
	}
	err := getBots(botsURL, botsAPIToken, &payload)
	if err != nil {
		return nil, err
	}
	return payload.Traders, nil
}

func (b *ResearchBot) GetWallet() (*wallet.VegaWallet, error) {
	var err error
	if b.wallet == nil {
		if b.WalletData.RecoveryPhrase == nil {
			return nil, fmt.Errorf("failed to get wallet for bot %s, recovery phrase is missing", b.Name)
		}
		b.wallet, err = wallet.GetVegaWalletSingleton(
			*b.WalletData.RecoveryPhrase, uint32(b.WalletData.Index),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get wallet for bot %s, %w", b.Name, err)
		}
	}
	return b.wallet, nil
}
