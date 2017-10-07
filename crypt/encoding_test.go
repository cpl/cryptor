package crypt_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	// Generate random data
	data := crypt.RandomData(32)
	// Encode then decode the data
	dData, err := crypt.Decode(crypt.Encode(data))
	if err != nil {
		t.Error(err)
	}

	// Compare decoded data with original data
	if !bytes.Equal(dData, data) {
		t.Error("data mismatch: initial data and encoded/decoded data")
	}
}

func TestStringEncodeDecode(t *testing.T) {
	t.Parallel()

	// Generate random data
	data := crypt.RandomData(32)
	// Encode data as string then decode the string as []byte
	dData, err := crypt.DecodeString(crypt.EncodeString(data))
	if err != nil {
		t.Error(err)
	}

	// Compare decoded string with original data
	if !bytes.Equal(dData, data) {
		t.Error("data mismatch: initial data and sencode/sdecode data")
	}
}

func TestDecodeErrors(t *testing.T) {
	t.Parallel()

	// Attempt decode invalid data
	if _, err := crypt.Decode([]byte{0, 1, 2, 3, 4}); err != nil {
		if err.Error() != "encoding/hex: odd length hex string" {
			t.Error("encoding/hex: decoded invalid data")
		}
	}

	// Attempt decode valid data
	if _, err := crypt.Decode([]byte{0, 1, 2, 3}); err != nil {
		if err.Error() != "encoding/hex: invalid byte: U+0000" {
			t.Error("encoding/hex: decoded invalid data")
		}
	}
}
