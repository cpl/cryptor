package hashing

import (
	"encoding/hex"
	"hash"

	"golang.org/x/crypto/blake2s"
)

// HashSize ...
const HashSize = blake2s.Size

// HashSum ...
type HashSum [HashSize]byte

// ToHex ...
func (h *HashSum) ToHex() string {
	return hex.EncodeToString(h[:])
}

// HashFunction ...
func HashFunction() hash.Hash {
	h, _ := blake2s.New256(nil)
	return h
}

// Hash ...
func Hash(sum *HashSum, data ...[]byte) {
	h := HashFunction()
	defer h.Reset()

	for _, set := range data {
		h.Write(set)
	}

	if sum == nil {
		return
	}

	h.Sum(sum[:0])
	return
}
