package crypt_test

import "testing"
import "github.com/thee-engineer/cryptor/crypt"

func BenchmarkEncryption(b *testing.B) {
	key := crypt.NewKey()
	dat := crypt.RandomData(16777216)
	b.ResetTimer()

	for count := 0; count < b.N; count++ {
		crypt.Encrypt(key, dat)
	}
}

func BenchmarkDecryption(b *testing.B) {
	key := crypt.NewKey()
	dat := crypt.RandomData(16777216)

	enc, err := crypt.Encrypt(key, dat)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()

	for count := 0; count < b.N; count++ {
		crypt.Decrypt(key, enc)
	}
}
