package hashing_test

import (
	"math/rand"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

func BenchmarkHash(b *testing.B) {
	rand.Seed(1)
	b.ResetTimer()

	for iteration := 0; iteration < b.N; iteration++ {
		hashing.Hash(crypt.RandomData(1024))
	}
}

func BenchmarkSum(b *testing.B) {
	rand.Seed(1)
	b.ResetTimer()

	for iteration := 0; iteration < b.N; iteration++ {
		hashing.Sum(crypt.RandomData(1024))
	}
}
