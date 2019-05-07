package hkdf_test

import (
	"crypto/hmac"
	"hash"
	"math/rand"
	"testing"

	"golang.org/x/crypto/blake2s"

	"cpl.li/go/cryptor/crypt"
)

func BenchmarkBLAKE(b *testing.B) {
	// generate BLAKE2s MAC
	rand.Seed(1)
	mac, err := blake2s.New256(crypt.RandomBytes(blake2s.Size))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	// write to MAC
	for iter := 0; iter < b.N; iter++ {
		mac.Write(crypt.RandomBytes(4096))
		mac.Sum(nil)
		mac.Reset()
	}
}

func localHash() hash.Hash {
	h, _ := blake2s.New256(nil)
	return h
}

func BenchmarkHMACwBLAKE(b *testing.B) {
	// generate BLAKE2s MAC
	rand.Seed(1)
	mac := hmac.New(localHash, crypt.RandomBytes(blake2s.Size))

	b.ResetTimer()

	// write to MAC
	for iter := 0; iter < b.N; iter++ {
		mac.Write(crypt.RandomBytes(4096))
		mac.Sum(nil)
		mac.Reset()
	}
}
