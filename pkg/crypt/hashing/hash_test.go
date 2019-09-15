package hashing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cpl.li/go/cryptor/pkg/crypt/hashing"
)

func assertHash(t *testing.T, data, expected string) {
	var sum hashing.Blake2sHash
	hashing.Hash(&sum, []byte(data))
	assert.Equal(t, expected, sum.ToHex(), "non matching hashes")
}

func TestHash(t *testing.T) {
	t.Parallel()

	assertHash(t, "",
		"69217a3079908094e11121d042354a7c1f55b6482ca1a51e1b250dfd1ed0eef9")
	assertHash(t, "Hello, World!",
		"ec9db904d636ef61f1421b2ba47112a4fa6b8964fd4a0a514834455c21df7812")
}

func TestHashNil(t *testing.T) {
	t.Parallel()

	data := []byte("Hello, World!")

	// request hash to be returned not passed in first argument
	hash0 := hashing.Hash(nil, data)
	assert.Equal(t, hash0.ToHex(),
		"ec9db904d636ef61f1421b2ba47112a4fa6b8964fd4a0a514834455c21df7812",
		"failed to match returned hash")

	// request hash to be passed in as argument
	var sum hashing.Blake2sHash
	hash1 := hashing.Hash(&sum, data)
	if hash1 != nil {
		t.Fatalf("hash function returned non-nil byte array %v\n", hash1[:])
	}

	if sum != *hash0 {
		t.Fatalf("failed to match hashes\n")
	}
}
