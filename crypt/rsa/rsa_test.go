package rsa_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/common/con"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/rsa"
)

func TestEncryption(t *testing.T) {
	t.Parallel()

	// Generate data and key
	data := crypt.RandomData(128)
	key := rsa.NewKey()

	// Encrypt data
	eData, err := rsa.Encrypt(&key.PublicKey, data)
	if err != nil {
		t.Error(err)
	}

	// Sanity check
	if bytes.Equal(eData, data) {
		t.Errorf("rsa encrypt | data still matches")
	}

	// Decrypt data
	dData, err := rsa.Decrypt(key, eData)
	if err != nil {
		t.Error(err)
	}

	// Check if data matches
	if !bytes.Equal(data, dData) {
		t.Errorf("rsa decrypt | data does not match")
	}
}

func TestEncryptionErrors(t *testing.T) {
	t.Parallel()

	// Generate data and key
	data := crypt.RandomData(con.KB)
	key := rsa.NewKey()

	_, err := rsa.Encrypt(&key.PublicKey, data)
	if err.Error() != "crypto/rsa: message too long for RSA public key size" {
		t.Error(err)
	}

	_, err = rsa.Decrypt(key, data)
	if err.Error() != "crypto/rsa: decryption error" {
		t.Error(err)
	}

	eData, err := rsa.Encrypt(&key.PublicKey, crypt.RandomData(10))
	if err != nil {
		t.Error(err)
	}
	nKey := rsa.NewKey()
	_, err = rsa.Decrypt(nKey, eData)
	if err.Error() != "crypto/rsa: decryption error" {
		t.Error(err)
	}

	eData[0] += 1
	_, err = rsa.Decrypt(key, eData)
	if err.Error() != "crypto/rsa: decryption error" {
		t.Error(err)
	}
}

func TestSignature(t *testing.T) {
	t.Parallel()

	data := crypt.RandomData(con.KB)
	key := rsa.NewKey()

	signature, err := rsa.Sign(key, data)
	if err != nil {
		t.Error(err)
	}

	if !rsa.Verify(&key.PublicKey, data, signature) {
		t.Errorf("rsa verify failed")
	}
}

func TestSignatureErrors(t *testing.T) {
	t.Parallel()

	data := crypt.RandomData(con.KB)
	key0 := rsa.NewKey()
	key1 := rsa.NewKey()

	signature, err := rsa.Sign(key0, data)
	if err != nil {
		t.Error()
	}

	if rsa.Verify(&key1.PublicKey, data, signature) {
		t.Errorf("rsa verified wrong key")
	}
	if rsa.Verify(&key0.PublicKey, data, data) {
		t.Errorf("rsa verified with wrong format")
	}
	if rsa.Verify(&key0.PublicKey, data, crypt.RandomData(512)) {
		t.Errorf("rsa verified random signature")
	}
	signature[0] += 1
	if rsa.Verify(&key0.PublicKey, data, signature) {
		t.Errorf("rsa verified bad signature")
	}
}
