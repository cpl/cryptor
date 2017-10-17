package crypt_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestSign(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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

func TestSignInvalid(t *testing.T) {
	t.Parallel()

	testData := crypt.RandomData(1024)

	key0, err := crypt.GenerateKey256()
	if err != nil {
		t.Error(err)
	}

	key1, err := crypt.GenerateKey256()
	if err != nil {
		t.Error(err)
	}

	sign, err := key0.Sign(testData)
	if err != nil {
		t.Error(err)
	}

	if key1.PublicKey.Verify(testData, sign) {
		t.Error("sign error: verified invalid signature")
	}

	if key1.PublicKey.Verify(crypt.RandomData(1024), sign) {
		t.Error("sign error: verified invalid signature")
	}

	if key1.PublicKey.Verify(testData, crypt.RandomData(64)) {
		t.Error("sign error: verified invalid signature")
	}

	if key1.PublicKey.Verify(testData, crypt.RandomData(16)) {
		t.Error("sign error: verified invalid signature")
	}

	if key1.PublicKey.Verify(testData, crypt.RandomData(100)) {
		t.Error("sign error: verified invalid signature")
	}
}

func TestSignSharedSecret(t *testing.T) {
	t.Parallel()

	key0, err := crypt.GenerateKey256()
	if err != nil {
		t.Error(err)
	}

	key1, err := crypt.GenerateKey256()
	if err != nil {
		t.Error(err)
	}

	sec, err := key0.GenerateShared(&key1.PublicKey, 16, 16)
	if err != nil {
		t.Error(err)
	}

	t.Log(len(key0.D.Bytes()))
	t.Log(len(sec))
}
