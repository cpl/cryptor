package hashing

import (
	"log"

	"golang.org/x/crypto/blake2b"
)

// Hash ...
func Hash(data []byte) []byte {
	h, err := blake2b.New512(nil)
	if err != nil {
		panic(err)
	}
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
	h, err := blake2b.New512(nil)
	if err != nil {
		panic(err)
	}
	w, err := h.Write(data)
	if err != nil {
		panic(err)
	}
	if w != len(data) {
		log.Panicf("blake2: write len %d does not match data len", w)
	}

	return h.Sum(data)
}
