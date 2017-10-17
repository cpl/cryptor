package crypt_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestSign(t *testing.T) {
	testData := crypt.RandomData(3000)

	key, err := crypt.GenerateKey256()
	if err != nil {
		t.Error(err)
	}

	sign, err := key.Sign(testData)
	if err != nil {
		t.Error(err)
	}

	if !key.PublicKey.Verify(testData, sign) {
		t.Error("sign error: failed to verify signature")
	}
}

func TestSignSmall(t *testing.T) {
	testData := crypt.RandomData(10)

	key, err := crypt.GenerateKey256()
	if err != nil {
		t.Error(err)
	}

	sign, err := key.Sign(testData)
	if err != nil {
		t.Error(err)
	}

	if !key.PublicKey.Verify(testData, sign) {
		t.Error("sign error: failed to verify signature")
	}
}

// TODO: Test for errors and invalid cases
