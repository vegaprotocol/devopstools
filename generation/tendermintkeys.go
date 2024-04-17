package generation

import (
	"encoding/base64"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

type TendermintKeys struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func GenerateTendermintKeys() TendermintKeys {
	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()
	address := pubKey.Address()
	return TendermintKeys{
		Address:    address.String(),
		PublicKey:  base64.StdEncoding.EncodeToString(pubKey.Bytes()),
		PrivateKey: base64.StdEncoding.EncodeToString(privKey.Bytes()),
	}
}
