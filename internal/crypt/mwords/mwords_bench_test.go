package mwords_test

import (
	"testing"

	"cpl.li/go/cryptor/internal/crypt"
	"cpl.li/go/cryptor/internal/crypt/mwords"
)

func BenchmarkEntropyToMnemonic(b *testing.B) {
	for iter := 0; iter < b.N; iter++ {
		_, err := mwords.EntropyToMnemonic(crypt.RandomBytes(32))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEntropyFromMnemonic(b *testing.B) {
	mnemonic, err := mwords.EntropyToMnemonic(crypt.RandomBytes(32))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for iter := 0; iter < b.N; iter++ {
		_, err := mwords.EntropyFromMnemonic(mnemonic)
		if err != nil {
			b.Fatal(err)
		}
	}
}
