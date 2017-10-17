// Package hashing provides SHA256 hashing algorithm. With plans for future
// hashing algorithms.
package hashing

import (
	"crypto/sha256"
	"hash"

	"github.com/thee-engineer/cryptor/crypt"
)

// SHA256 returns the SHA256 hash of the given data.
func SHA256(data []byte) hash.Hash {
	h := sha256.New()
	if _, err := h.Write(data); err != nil {
		panic(err)
	}

	return h
}

// SHA256Digest returns the result of the hash function as []byte array.
func SHA256Digest(data []byte) []byte {
	return SHA256(data).Sum(nil)
}

// SHA256HexDigest returns the result of the hash function
// as a hex encoded string.
func SHA256HexDigest(data []byte) string {
	return crypt.EncodeString(SHA256Digest(data))
}
