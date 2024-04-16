package ethutils

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/params"
)

func EtherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func TokenToFullTokens(amount *big.Int, decimals uint8) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(amount), big.NewFloat(math.Pow10(int(decimals))))
}

func TimestampNonce() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

func VegaPubKeyToByte32(pubKey string) (byte32PubKey [32]byte, err error) {
	var bytePubKey []byte
	bytePubKey, err = hex.DecodeString(pubKey)
	if err != nil {
		return
	}
	copy(byte32PubKey[:], bytePubKey)
	return
}
