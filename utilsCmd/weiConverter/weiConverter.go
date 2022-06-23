package weiConverter

import (
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"strings"
)

func WeiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

func WeiToGwei(wei *big.Int) *big.Float {
	//return new(big.Int).Div(wei, big.NewInt(params.GWei))
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.GWei))
	//return wei
}

func GweiTowei(gwei *big.Float) *big.Int {
	//return new(big.Int).Mul(gwei, big.NewInt(params.GWei))
	truncInt, _ := gwei.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.GWei))
	fracStr := strings.Split(fmt.Sprintf("%.18f", gwei), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func EtherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

// ParseBigFloat parse string value to big.Float
func ParseBigFloat(value string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	return f, err
}
