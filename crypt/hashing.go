package crypt

import (
	"crypto/sha256"
	"hash"
	"io"
	"os"

	"golang.org/x/crypto/sha3"
)

// SHA256File takes a file path and returns the hash for the file content
func SHA256File(path string) (hash.Hash, error) {
	// Open file for reading
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create new SHA256 Hash
	h := sha256.New()
	// Hash file content
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}

	return h, nil
}

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
