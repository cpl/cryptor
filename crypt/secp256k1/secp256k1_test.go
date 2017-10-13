package secp256k1_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/secp256k1"
)

func TestSign(t *testing.T) {
	msg := crypt.RandomData(32)
	key0, err := crypt.GenerateKey256()
	if err != nil {
		t.Error(err)
	}

	key1, err := crypt.GenerateKey256()
	if err != nil {
		t.Error(err)
	}

	secret, err := key1.GenerateShared(&key0.PublicKey, 16, 16)
	if err != nil {
		t.Error(err)
	}

	_, err = secp256k1.Sign(msg, secret)
	if err != nil {
		t.Error(err)
	}
}
