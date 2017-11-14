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

	aes.NewKeyFromBytes(key.Bytes())
	aes.NewKeyFromBytes(crypt.RandomData(aes.KeySize))
	aes.NewKeyFromString(key.String())
	aes.NewKeyFromString("")

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
		"64a868d4b23af696d3734d0b814d04cdd1ac280128e97653a05f32b49c13a29a",
		"65e4e5bf4cb95801ead6d6dfb3367bfcd3c7cfebd62efaee312f2b6d40ff4f5c",
		"00c4ce73053141d6b53016dc02bc4efaa65f38ed457ed1f8b28eb65499b7f555",
		"bae68a56d1acea882202bcc3fe8857235bb410df38c9ae28e7f2e44100f9aced",
		"cd65b591e9d0c68359f344c9c693e75b51544aa4ba88b32d7555517c447427a2"}

	// Test all values
	for index, test := range testValues {
		if testResult[index] != string(aes.NewKeyFromPassword(test).Encode()) {
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
	if err != nil {
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
	if !bytes.Equal(key.Bytes(), make([]byte, aes.KeySize)) {
		t.Error("aes key: zeroing failed")
	}
}
