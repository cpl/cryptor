package hashing

import (
	"hash"
	"log"

	"golang.org/x/crypto/blake2b"
)

// HashFunction ...
func HashFunction() hash.Hash {
	if h, err := blake2b.New256(nil); err != nil {
		panic(err)
	} else {
		return h
	}
}

// Hash ...
func Hash(data []byte) []byte {
	h := HashFunction()
	w, err := h.Write(data)
	if err != nil {
		panic(err)
	}
	if w != len(data) {
		log.Panicf("blake2: write len %d does not match data len", w)
	}

	return h.Sum(nil)
}

// Sum ...
func Sum(data []byte) []byte {
	h := HashFunction()
	w, err := h.Write(data)
	if err != nil {
		panic(err)
	}
	if w != len(data) {
		log.Panicf("blake2: write len %d does not match data len", w)
	}

	return h.Sum(data)
}
