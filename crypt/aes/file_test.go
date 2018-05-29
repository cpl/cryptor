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
	initialData, err := ioutil.ReadFile("../../test/testfile.txt")
	if err != nil {
		t.Error(err)
	}

	// Encrypt files
	if err := aes.EncryptFiles("testpassword",
		"../../test/testfile.txt", "../../test/testfile.log"); err != nil {
		t.Error(err)
	}

	// Fail decryption
	if err := aes.DecryptFiles("wrong password",
		"../../test/testfile.txt", "../../test/testfile.log"); err == nil {
		t.Error("decrypted with wrong password")
	}

	// Decrypt files
	if err := aes.DecryptFiles("testpassword",
		"../../test/testfile.txt", "../../test/testfile.log"); err != nil {
		t.Error(err)
	}

	// Read final data
	finalData, err := ioutil.ReadFile("../../test/testfile.txt")
	if err != nil {
		t.Error(err)
	}

	// Compare
	if !bytes.Equal(initialData, finalData) {
		t.Errorf("aes file encryption, files do not match")
	}
}
