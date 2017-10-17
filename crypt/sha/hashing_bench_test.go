package crypt_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func BenchmarkSHA256(b *testing.B) {
	for count := 0; count < b.N; count++ {
		crypt.SHA256Data(crypt.RandomData(1048576))
	}
}

func BenchmarkSHA512(b *testing.B) {
	for count := 0; count < b.N; count++ {
		crypt.SHA512Data(crypt.RandomData(1048576))
	}
}
