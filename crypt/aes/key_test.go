package aes_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

func TestKey(t *testing.T) {
	t.Parallel()

	key := aes.NewKey()

	_, err := aes.NewKeyFromBytes(key.Bytes())
	if err != nil {
		t.Error(err)
	}

	_, err = aes.NewKeyFromBytes(crypt.RandomData(crypt.KeySize))
	if err != nil {
		t.Error(err)
	}

	_, err = aes.NewKeyFromString(key.String())
	if err != nil {
		t.Error(err)
	}

	_, err = aes.NewKeyFromString("")
	if err == nil {
		t.Error("aes key: created key from empty string")
	}

	_, err = aes.NewKeyFromBytes(crypt.RandomData(20))
	if err == nil {
		t.Error("aes key: create key from 20 bytes")
	}

	key.Encode()
}

func TestKeyFromPassword(t *testing.T) {
	t.Parallel()

	// Test passwords
	var testValues = []string{
		"",
		"0123",
		"testinglargerpassword",
		"5rdZ<Q/{Uf6@Ed!~uF8T(9Ad8S<+{w9,",
		">~Ph]p`}$_QSW6#wrfo$/a<fH+PLz,Ycv;tY#Y4b&2uNwC7BKT~GZ<HsPXt(}"}

	// Expected test results for password key derivation
	var testResult = []string{
		"5537f231b6f9f0fcf5a03c5043ced5f2d8b8f69de6ae46f7de5f269b38f5161a",
		"414885cd49145a4a6656ede93728f2a005899c9891c5c42839784903d7e99f6c",
		"33c54d2ad9fbb0d9745a927b515d4872d52507e884256aa45b5d6a43f6197bfb",
		"d4e55437bae6bed6a59ca25a9647d74ebf98772d6eafc88394baf3c0b790d973",
		"d380bef98b6d843d767863021db1106d13f36fd54b34f4064ca9476bd3191c3a"}

	// Test all values
	for index, test := range testValues {
		key, err := aes.NewKeyFromPassword(test)
		if err != nil {
			t.Error(err)
		}

		if testResult[index] != string(key.Encode()) {
			t.Logf("expected: %s\n", testResult[index])
			t.Logf("got: %s\n", key.String())

			t.Error("aes key error: wrong key derivation")
		}
	}
}

func TestNewKeyFromStringError(t *testing.T) {
	t.Parallel()

	// Test invalid string
	_, err := aes.NewKeyFromString(string(crypt.RandomData(64)))
	if err == nil {
		t.Error("aes key error: invalid hex string used as key")
	}

	// Test large string
	_, err = aes.NewKeyFromString(string(crypt.RandomData(100)))
	if err == nil {
		t.Error("aes key error: invalid hex string used as key")
	}

	// Test small string
	_, err = aes.NewKeyFromString(string(crypt.RandomData(32)))
	if err == nil {
		t.Error("aes key error: invalid hex string used as key")
	}

	// Test string of white spaces
	key, err := aes.NewKeyFromString(
		"\n    \t     \n\n   \n \t        \n     \n\t\n     \t     \n      ")
	if err == nil {
		t.Error("aes key error: invalid hex string used as key")
	}

	if key != aes.NullKey {
		t.Error("aes key error: obtained invalid aes key from empty hex")
	}
}

func TestKeyZeroing(t *testing.T) {
	t.Parallel()

	// Generate AES key
	key := aes.NewKey()
	// Zero aes key
	crypt.ZeroBytes(key[:])

	// Check that zeroing succeeded
	if !bytes.Equal(key.Bytes(), make([]byte, crypt.KeySize)) {
		t.Error("aes key: zeroing failed")
	}
}
