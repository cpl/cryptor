package hashing_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/crypt/hashing"
)

func TestHashing(t *testing.T) {
	t.Parallel()

	// Use hash of "Hello, World!" generated outside cryptor
	eHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"

	// Hash "Hello, World!" using cryptor/crypt
	hash := hashing.SHA256HexDigest([]byte("Hello, World!"))

	// Compare hashes
	if eHash != hash {
		t.Error("sha256 error: hashes don't match")
	}
}
