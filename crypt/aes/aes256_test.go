package aes_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

func TestAESEncryption(t *testing.T) {
	t.Parallel()

	// Generate key from string
	key, err := aes.NewKeyFromString(
		"0873eacc863d4748b237fd4d4c877926aa111092c14e19d9f5730479c7fb92a6")
	if err != nil {
		t.Error(err)
	}

	msg := []byte("Hello, World!")

	// Encrypt msg
	eMsg, err := aes.Encrypt(key, msg)
	if err != nil {
		t.Error(err)
	}

	// Decrypt msg
	dMsg, err := aes.Decrypt(key, eMsg)
	if err != nil {
		t.Error(err)
	}

	// Compare initial msg with decrypted msg
	if !bytes.Equal(dMsg, msg) {
		t.Error("data mismatch: initial msg and encrypted->decrypted msg")
	}
}

func TestCryptoErrors(t *testing.T) {
	t.Parallel()

	// Generate key from string
	key, err := aes.NewKeyFromString(
		"0873eacc863d4748b237fd4d4c877926aa111092c14e19d9f5730479c7fb92a6")
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("Hello, World!")

	// Encrypt msg with null key
	eMsg, err := aes.Encrypt(aes.NullKey, msg)
	if err != nil {
		t.Error(err)
	}

	// Decrypt with key, msg encrypted with null key
	_, err = aes.Decrypt(key, eMsg)
	if err.Error() != "cipher: message authentication failed" {
		t.Error(err)
	}

	// Decrypt random data which is smaller than the nonce
	_, err = aes.Decrypt(key, crypt.RandomData(10))
	if err.Error() != "invalid nonce" {
		t.Error(err)
	}

	// Add random data to valid encrypted msg
	mMsg := append(crypt.RandomData(12), eMsg[12:]...)

	// Attempt to decrypt invalid msg
	_, err = aes.Decrypt(aes.NullKey, mMsg)
	if err.Error() != "cipher: message authentication failed" {
		t.Error(err)
	}
}

func TestEncryptDecryptNullKey(t *testing.T) {
	t.Parallel()
	// Random data
	data := crypt.RandomData(100)

	// Encrypt
	eData, err := aes.Encrypt(aes.NullKey, data)
	if err != nil {
		t.Error(err)
	}

	// Decrypt
	dData, err := aes.Decrypt(aes.NullKey, eData)
	if err != nil {
		t.Error(err)
	}

	// Compare
	if !bytes.Equal(data, dData) {
		t.Error("data mismatch: initial msg and encrypted->decrypted msg")
	}
}
