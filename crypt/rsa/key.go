package rsa

import (
	"crypto/rand"
	"crypto/rsa"
)

// KeySizeBits ...
const KeySizeBits = 4096

// NewKey ...
func NewKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, KeySizeBits)
	if err != nil {
		panic(err)
	}
	if err := key.Validate(); err != nil {
		panic(err)
	}
	key.Precompute()
	return key
}
