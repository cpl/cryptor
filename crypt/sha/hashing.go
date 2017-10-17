package crypt

import (
	"crypto/sha256"
	"hash"

	"golang.org/x/crypto/sha3"
)

// SHA256Data returns the SHA256 hash of the given data
func SHA256Data(data []byte) hash.Hash {
	h := sha256.New()
	if _, err := h.Write(data); err != nil {
		panic(err)
	}

	return h
}

// SHA512Data takes multiple []byte arrays and returns the combined hash
func SHA512Data(data []byte) hash.Hash {
	h := sha3.New512()
	if _, err := h.Write(data); err != nil {
		panic(err)
	}
	return h
}
