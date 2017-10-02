package crypt

import (
	"testing"
)

func TestHashing(t *testing.T) {
	t.Parallel()

	// Use hash of "Hello, World!" generated outside cryptor
	eHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
	// Hash "Hello, World!" using cryptor/crypt
	hash := SHA256Data([]byte("Hello, World!"))
	// Encode the hash as a string
	sHash := EncodeString(hash.Sum(nil))

	// Compare hashes
	if eHash != sHash {
		t.Error("data mismatch: hash and hash")
	}

	// Attempt to hash source file
	_, err := SHA256File("hashing.go")
	if err != nil {
		t.Error(err)
	}
}
