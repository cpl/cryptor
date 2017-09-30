package crypt

import (
	"bytes"
	"testing"
)

func TestCrypt(t *testing.T) {
	key := NewKeyFromString("0873eacc863d4748b237fd4d4c877926aa111092c14e19d9f5730479c7fb92a6")
	msg := []byte("Hello, World!")

	eMsg, err := Encrypt(key, msg)
	if err != nil {
		t.Error(err)
	}

	dMsg, err := Decrypt(key, eMsg)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(dMsg, msg) != 0 {
		t.Error("Initial message does not match encrypted->decrypted msg!")
	}
}

func TestCryptoErrors(t *testing.T) {
	key := NewKeyFromString("0873eacc863d4748b237fd4d4c877926aa111092c14e19d9f5730479c7fb92a6")
	msg := []byte("Hello, World!")

	eMsg, err := Encrypt(NullKey, msg)
	if err != nil {
		t.Error(err)
	}

	_, err = Decrypt(key, eMsg)
	if err.Error() != "cipher: message authentication failed" {
		t.Error(err)
	}

	_, err = Decrypt(key, RandomData(10))
	if err.Error() != "Invalid nonce" {
		t.Error(err)
	}

	mMsg := append(RandomData(12), eMsg[12:]...)
	_, err = Decrypt(NullKey, mMsg)
	if err.Error() != "cipher: message authentication failed" {
		t.Error(err)
	}
}
