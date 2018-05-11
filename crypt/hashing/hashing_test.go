package hashing_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt/hashing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestHashAndSum(t *testing.T) {
	t.Parallel()
	dataSets := make([][]byte, 10)
	hashes := make([][]byte, 10)
	sums := make([][]byte, 10)

	for i, _ := range dataSets {
		dataSets[i] = crypt.RandomData(1024)
		hashes[i] = hashing.Hash(dataSets[i])
		sums[i] = hashing.Sum(dataSets[i])
	}

	for i, _ := range dataSets {
		if !bytes.Equal(hashes[i], hashing.Hash(dataSets[i])) {
			t.Errorf("failed hash")
		}
		if !bytes.Equal(sums[i], hashing.Sum(dataSets[i])) {
			t.Errorf("failed checksum")
		}
	}
}
