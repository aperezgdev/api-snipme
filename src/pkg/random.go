package pkg

import (
	"crypto/rand"
	"io"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RandIntFunc func(reader io.Reader, max *big.Int) (*big.Int, error)

func GenerateShortCodeWithRand(n int, randFunc RandIntFunc) (string, error) {
	code := make([]byte, n)
	for i := range code {
		num, err := randFunc(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}
	return string(code), nil
}

func GenerateShortCode(n int) (string, error) {
	return GenerateShortCodeWithRand(n, rand.Int)
}
