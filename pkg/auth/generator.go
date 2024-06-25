package auth

import (
	"crypto/rand"
	"math/big"
)

const (
	maxNumber int64 = 9999
	minNumber int64 = 1000
)

type defaultGenerator struct {
}

func (g *defaultGenerator) generateCode() LoginCode {
	randomNumber, _ := rand.Int(rand.Reader, big.NewInt(maxNumber-minNumber))
	return LoginCode(minNumber + randomNumber.Int64())
}
