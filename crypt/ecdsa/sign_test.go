package crypt_test

import (
	"bytes"
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

func TestSignVerifySharedSecret(t *testing.T) {
	t.Parallel()

	prv0, prv1, err := generateKeyPair256()
	if err != nil {
		t.Error(err)
	}

	sec0, err := prv0.GenerateShared(&prv1.PublicKey, 16, 16)
	if err != nil {
		t.Error(err)
	}

	sec1, err := prv1.GenerateShared(&prv0.PublicKey, 16, 16)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(sec0, sec1) {
		t.Error("shared keys: mismatched secrets")
	}

	sec1Key, err := crypt.Key256FromSecret(sec1)
	if err != nil {
		t.Error(err)
	}

	sec0Key, err := crypt.Key256FromSecret(sec0)
	if err != nil {
		t.Error(err)
	}

	if !compareKeys(sec0Key, sec1Key) {
		t.Error("shared keys: not matching")
	}

	testData := crypt.RandomData(1024)

	sign0, err := sec0Key.Sign(testData)
	if err != nil {
		t.Error(err)
	}

	sign1, err := sec1Key.Sign(testData)
	if err != nil {
		t.Error(err)
	}

	if !sec0Key.Verify(testData, sign0) {
		t.Error("sign: failed to verify test data with signature")
	}

	if !sec0Key.Verify(testData, sign1) {
		t.Error("sign: failed to verify test data with signature")
	}

	if !sec1Key.Verify(testData, sign0) {
		t.Error("sign: failed to verify test data with signature")
	}

	if !sec1Key.Verify(testData, sign1) {
		t.Error("sign: failed to verify test data with signature")
	}

	// Invalid signature
	if sec0Key.Verify(testData, crypt.RandomData(32)) {
		t.Error("sign: verified invalid signature")
	}

	// Invalid data
	if sec0Key.Verify(crypt.RandomData(1024), sign0) {
		t.Error("sign: verified invalid data")
	}
}
