package aes_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/thee-engineer/cryptor/crypt/aes"
)

func TestFileEncryption(t *testing.T) {
	t.Parallel()

	// Memorize intial test data
	initialData, err := ioutil.ReadFile("testfile.txt")
	if err != nil {
		t.Error(err)
	}

	// Encrypt files
	if err := aes.EncryptFiles("testpassword",
		"testfile.txt", "testfile.log"); err != nil {
		t.Error(err)
	}

	// Fail decryption
	if err := aes.DecryptFiles("wrong password",
		"testfile.txt", "testfile.log"); err == nil {
		t.Error("decrypted with wrong password")
	}

	// Decrypt files
	if err := aes.DecryptFiles("testpassword",
		"testfile.txt", "testfile.log"); err != nil {
		t.Error(err)
	}

	// Read final data
	finalData, err := ioutil.ReadFile("testfile.txt")
	if err != nil {
		t.Error(err)
	}

	// Compare
	if !bytes.Equal(initialData, finalData) {
		t.Errorf("aes file encryption, files do not match")
	}

}
