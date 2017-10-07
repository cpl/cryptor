package crypt_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestHashing(t *testing.T) {
	t.Parallel()

	// Use hash of "Hello, World!" generated outside cryptor
	eHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
	// Hash "Hello, World!" using cryptor/crypt
	hash := crypt.SHA256Data([]byte("Hello, World!"))
	// Encode the hash as a string
	sHash := crypt.EncodeString(hash.Sum(nil))

	// Compare hashes
	if eHash != sHash {
		t.Error("data mismatch: hash and hash")
	}

	// Attempt to hash source file
	_, err := crypt.SHA256File("hashing.go")
	if err != nil {
		t.Error(err)
	}
}

// Test SHA512
func TestSHA512(t *testing.T) {
	t.Parallel()

	// Test data
	testData := []string{" ", "hello", "µcryptor", "\x00"}
	results := []string{
		"e307daea2f0168daa1318e2faa2d67791e9d8e03692a6f7d1eb974e664fe721e81a47b4cf3d0eb19ae5d57afa19a095941cad5a5c050774ad56a8e5e21105757",
		"75d527c368f2efe848ecf6b073a36767800805e9eef2b1857d5f984f036eb6df891d75f72d9b154518c1cd58835286d1da9a38deba3de98b5a53e5ed78a84976",
		"eb927a04901cc128cee0b5cd0a931f2326a362ecd3f7c0b5ed2ab25c5b454efbb56633238e1c384daed806fab5c2511a9fb78a23bb4f965afcaf2189aa89c3bb",
		"7127aab211f82a18d06cf7578ff49d5089017944139aa60d8bee057811a15fb55a53887600a3eceba004de51105139f32506fe5b53e1913bfa6b32e716fe97da"}

	// Hash each test and check for expected results
	for index, data := range testData {
		hash := crypt.SHA512Data([]byte(data))
		if crypt.EncodeString(hash.Sum(nil)) != results[index] {
			t.Error("hash error: wrong hash result")
		}
	}
}
