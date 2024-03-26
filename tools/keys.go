package tools

import "encoding/hex"

func HexKeyToByte32(hexKey string) ([32]byte, error) {
	bytePubKey, err := hex.DecodeString(hexKey)
	if err != nil {
		return [32]byte{}, err
	}

	var byte32PubKey [32]byte
	copy(byte32PubKey[:], bytePubKey)

	return byte32PubKey, nil
}
