package generation

import (
	"crypto/rand"
	"math/big"
)

func Password() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-!@#$%^&*()_+="
	lettersLen := big.NewInt(int64(len(letters)))
	// randomise password len in range 51-65
	passwordLenB, err := rand.Int(rand.Reader, big.NewInt(15))
	if err != nil {
		return "", err
	}
	passwordLen := 51 + passwordLenB.Int64()
	ret := make([]byte, passwordLen)
	for i := int64(0); i < passwordLen; i++ {
		num, err := rand.Int(rand.Reader, lettersLen)
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
