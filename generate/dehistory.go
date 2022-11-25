package generate

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"

	"github.com/ipfs/kubo/config"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

func GenerateDeHistoryIdentity(seed string) (config.Identity, error) {
	ident := config.Identity{}

	var sk crypto.PrivKey
	var pk crypto.PubKey

	// Everything > 32 bytes is ignored in GenerateEd25519Key so do a little pre hashing
	seedHash := sha256.Sum256([]byte(seed))

	priv, pub, err := crypto.GenerateEd25519Key(bytes.NewReader(seedHash[:]))
	if err != nil {
		return ident, err
	}

	sk = priv
	pk = pub

	skbytes, err := crypto.MarshalPrivateKey(sk)
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	return ident, nil
}
