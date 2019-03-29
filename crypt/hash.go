package crypt

import (
	"encoding/hex"
	"hash"

	"golang.org/x/crypto/blake2s"
)

// Blake2sHash represents a BLAKE2s 256 checksum of 32 bytes.
type Blake2sHash [blake2s.Size]byte

// ToHex exports the 32 byte checksum to a hex string.
func (h Blake2sHash) ToHex() string {
	return hex.EncodeToString(h[:])
}

// HashFunction returns the blake2s hash.Hash, by dropping the error for
// invalid key size as defined in `newDigest` function:
// https://github.com/golang/crypto/blob/master/blake2s/blake2s.go
func HashFunction() hash.Hash {
	h, _ := blake2s.New256(nil)
	return h
}

// Hash can take multiple byte arrays as input and compute their sum using
// BLAKE2s 256 hashing. The checksum is returned in the first argument as
// it must be a pointer to a byte array.
func Hash(sum *Blake2sHash, data ...[]byte) {
	// create hash
	hash := HashFunction()
	defer hash.Reset()

	// iterate byte arrays
	for _, set := range data {
		hash.Write(set)
	}

	hash.Sum(sum[:0])
}
