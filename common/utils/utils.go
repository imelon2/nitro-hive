package utils

import (
	"math/big"
	"strings"
)

func Unhexlify(input string) string {
	if strings.HasPrefix(input, "0x") {
		return strings.TrimPrefix(input, "0x")
	}
	return input
}

func FillBigIntArray(size int, value *big.Int) []*big.Int {
	arr := make([]*big.Int, size)
	for i := range arr {
		arr[i] = value
	}
	return arr
}
