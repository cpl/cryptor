package hashing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cpl.li/go/cryptor/internal/crypt/hashing"
)

func assertHash(t *testing.T, data, expected string) {
	var sum hashing.HashSum
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
