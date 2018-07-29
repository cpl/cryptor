package hashing_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/hashing"
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

func TestHashLength(t *testing.T) {
	t.Parallel()
	hashLen := len(hashing.Hash(crypt.RandomData(1)))
	if hashLen != hashing.HashSize {
		t.Fatalf("invalid hash size, expected %d got %d",
			hashing.HashSize, hashLen)
	}
}

func TestHashMultiple(t *testing.T) {
	t.Parallel()

	data0 := crypt.RandomData(100)
	data1 := crypt.RandomData(100)

	hashAll := hashing.Hash(data0, data1)
	hash0 := hashing.Hash(data0)
	hash1 := hashing.Hash(data1)

	if bytes.Equal(hashAll, hash0) || bytes.Equal(hashAll, hash1) {
		t.Errorf("failed hash, multi-hash matches single hash")
	}
}
