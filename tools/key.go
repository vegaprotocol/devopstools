package tools

import "encoding/hex"

func KeyAsByte32(pubKey string) (byte32PubKey [32]byte, err error) {
	var bytePubKey []byte
	bytePubKey, err = hex.DecodeString(pubKey)
	if err != nil {
		return
	}
	copy(byte32PubKey[:], bytePubKey)
	return
}
