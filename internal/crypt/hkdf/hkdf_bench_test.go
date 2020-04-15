package hkdf_test

import (
	"crypto/hmac"
	"hash"
	"math/rand"
	"testing"

	"golang.org/x/crypto/blake2s"

	"cpl.li/go/cryptor/internal/crypt"
)

func BenchmarkBLAKE(b *testing.B) {
	rand.Seed(1)
	mac, err := blake2s.New256(crypt.RandomBytes(blake2s.Size))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for iter := 0; iter < b.N; iter++ {
		data := crypt.RandomBytes(4096)
		mac.Write(data)
		mac.Sum(nil)
		mac.Reset()
	}
}

func localHash() hash.Hash {
	h, _ := blake2s.New256(nil)
	return h
}

func BenchmarkHMACwBLAKE(b *testing.B) {
	rand.Seed(1)
	mac := hmac.New(localHash, crypt.RandomBytes(blake2s.Size))

	b.ResetTimer()

	for iter := 0; iter < b.N; iter++ {
		data := crypt.RandomBytes(4096)
		mac.Write(data)
		mac.Sum(nil)
		mac.Reset()
	}
}
