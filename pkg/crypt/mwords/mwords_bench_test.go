package mwords_test

import (
	"testing"

	"cpl.li/go/cryptor/pkg/crypt"
	"cpl.li/go/cryptor/pkg/crypt/mwords"
)

func BenchmarkEntropyToMnemonic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mwords.EntropyToMnemonic(crypt.RandomBytes(32))
	}
}

func BenchmarkEntropyFromMnemonic(b *testing.B) {
	mnemonic, err := mwords.EntropyToMnemonic(crypt.RandomBytes(32))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mwords.EntropyFromMnemonic(mnemonic)
	}
}
