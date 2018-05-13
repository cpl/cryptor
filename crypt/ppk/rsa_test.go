package ppk_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/common/con"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/ppk"
)

func TestEncryption(t *testing.T) {
	t.Parallel()

	// Generate data and key
	data := crypt.RandomData(128)
	key := ppk.NewKey()

	// Encrypt data
	eData, err := ppk.Encrypt(&key.PublicKey, data)
	if err != nil {
		t.Error(err)
	}

	// Sanity check
	if bytes.Equal(eData, data) {
		t.Errorf("ppk encrypt | data still matches")
	}

	// Decrypt data
	dData, err := ppk.Decrypt(key, eData)
	if err != nil {
		t.Error(err)
	}

	// Check if data matches
	if !bytes.Equal(data, dData) {
		t.Errorf("ppk decrypt | data does not match")
	}
}

func TestEncryptionErrors(t *testing.T) {
	t.Parallel()

	// Generate data and key
	data := crypt.RandomData(con.KB)
	key := ppk.NewKey()

	_, err := ppk.Encrypt(&key.PublicKey, data)
	if err.Error() != "crypto/rsa: message too long for RSA public key size" {
		t.Error(err)
	}

	_, err = ppk.Decrypt(key, data)
	if err.Error() != "crypto/rsa: decryption error" {
		t.Error(err)
	}

	eData, err := ppk.Encrypt(&key.PublicKey, crypt.RandomData(10))
	if err != nil {
		t.Error(err)
	}
	nKey := ppk.NewKey()
	_, err = ppk.Decrypt(nKey, eData)
	if err.Error() != "crypto/rsa: decryption error" {
		t.Error(err)
	}

	eData[0] += 1
	_, err = ppk.Decrypt(key, eData)
	if err.Error() != "crypto/rsa: decryption error" {
		t.Error(err)
	}
}

func TestSignature(t *testing.T) {
	t.Parallel()

	data := crypt.RandomData(con.KB)
	key := ppk.NewKey()

	signature, err := ppk.Sign(key, data)
	if err != nil {
		t.Error(err)
	}

	if !ppk.Verify(&key.PublicKey, data, signature) {
		t.Errorf("ppk verify failed")
	}
}

func TestSignatureErrors(t *testing.T) {
	t.Parallel()

	data := crypt.RandomData(con.KB)
	key0 := ppk.NewKey()
	key1 := ppk.NewKey()

	signature, err := ppk.Sign(key0, data)
	if err != nil {
		t.Error()
	}

	if ppk.Verify(&key1.PublicKey, data, signature) {
		t.Errorf("ppk verified wrong key")
	}
	if ppk.Verify(&key0.PublicKey, data, data) {
		t.Errorf("ppk verified with wrong format")
	}
	if ppk.Verify(&key0.PublicKey, data, crypt.RandomData(512)) {
		t.Errorf("ppk verified random signature")
	}
	signature[0] += 1
	if ppk.Verify(&key0.PublicKey, data, signature) {
		t.Errorf("ppk verified bad signature")
	}
}
