package aes_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/scrypt"
)

func TestFileEncryption(t *testing.T) {
	t.Parallel()

	// Generate key using scrypt
	keyBytes := scrypt.Scrypt("testpassword", []byte("salt"))
	key, err := aes.NewKeyFromBytes(keyBytes)
	if err != nil {
		t.Error(err)
	}

	initialData, err := ioutil.ReadFile("testfile.txt")
	if err != nil {
		t.Error(err)
	}

	// Encrypt files
	if err := aes.EncryptFiles(key, "testfile.txt", "testfile.log"); err != nil {
		t.Error(err)
	}

	// Decrypt files
	if err := aes.DecryptFiles(key, "testfile.txt", "testfile.log"); err != nil {
		t.Error(err)
	}

	finalData, err := ioutil.ReadFile("testfile.txt")
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(initialData, finalData) {
		t.Errorf("aes file encryption, files do not match")
	}
}
