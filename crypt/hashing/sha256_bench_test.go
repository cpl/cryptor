package hashing_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

func BenchmarkSHA256(b *testing.B) {
	for count := 0; count < b.N; count++ {
		hashing.SHA256(crypt.RandomData(1048576))
	}
}

func BenchmarkSHA256HexDigest(b *testing.B) {
	for count := 0; count < b.N; count++ {
		hashing.SHA256Digest(crypt.RandomData(1048576))
	}
}
