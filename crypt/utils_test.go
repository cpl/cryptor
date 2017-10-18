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

	// Display all zeros (hopefully)
	t.Log(testData)
}

func TestZeroBytesMultiple(t *testing.T) {
	t.Parallel()

	// Generate three sets of test data
	testData0 := crypt.RandomData(100)
	testData1 := crypt.RandomData(5000)
	testData2 := crypt.RandomData(10000)

	// Zero all data
	crypt.ZeroBytes(testData0, testData1, testData2)

	// Check for zeros
	if !bytes.Equal(testData0, make([]byte, 100)) {
		t.Error("zero bytes: failed to zero all bytes")
	}
	if !bytes.Equal(testData1, make([]byte, 5000)) {
		t.Error("zero bytes: failed to zero all bytes")
	}
	if !bytes.Equal(testData2, make([]byte, 10000)) {
		t.Error("zero bytes: failed to zero all bytes")
	}
}
