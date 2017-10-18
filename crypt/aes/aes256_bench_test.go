package aes_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

func BenchmarkAESEncryption(b *testing.B) {
	key := aes.NewKey()
	dat := crypt.RandomData(16777216)
	b.ResetTimer()

	for count := 0; count < b.N; count++ {
		aes.Encrypt(key, dat)
	}
}

func BenchmarkAESDecryption(b *testing.B) {
	key := aes.NewKey()
	dat := crypt.RandomData(16777216)

	enc, err := aes.Encrypt(key, dat)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()

	for count := 0; count < b.N; count++ {
		aes.Decrypt(key, enc)
	}
}
