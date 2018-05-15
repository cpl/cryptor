package b64_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/encode/b64"
	"github.com/thee-engineer/cryptor/utils"
)

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	// Generate random data
	data := crypt.RandomData(32)
	// Encode then decode the data
	dData, err := b64.Decode(b64.Encode(data))
	utils.CheckErrTest(err, t)

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
	dData, err := b64.DecodeString(b64.EncodeString(data))
	utils.CheckErrTest(err, t)

	// Compare decoded string with original data
	if !bytes.Equal(dData, data) {
		t.Error("data mismatch: initial data and sencode/sdecode data")
	}
}

func TestDecodeErrors(t *testing.T) {
	t.Parallel()

	// Attempt decode invalid data
	if _, err := b64.Decode([]byte{0, 1, 2, 3, 4}); err != nil {
		if err == nil {
			t.Error("encoding/base64: decoded invalid data")
		}
	}

	// Attempt decode valid data
	if _, err := b64.Decode([]byte{0, 1, 2, 3}); err != nil {
		if err.Error() != "illegal base64 data at input byte 0" {
			t.Error(err)
			t.Error("encoding/base64: decoded invalid data")
		}
	}
}
