package crypt_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestZeroBytes(t *testing.T) {
	t.Parallel()

	testData := crypt.RandomData(1024)
	zeroBytes := make([]byte, 1024)

	crypt.ZeroBytes(testData)

	if len(testData) != 1024 {
		t.Error("zero bytes: detected change in data size")
	}

	if !bytes.Equal(zeroBytes, testData) {
		t.Error("zero bytes: failed to zero all bytes")
	}

	t.Log(testData)
}
