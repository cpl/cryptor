package crypt

import (
	"testing"
)

func TestHashing(t *testing.T) {
	eHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
	hash := SHA256Data([]byte("Hello, World!"))
	sHash := EncodeString(hash.Sum(nil))

	if eHash != sHash {
		t.Error("Hash did not match expected hash")
	}

	_, err := SHA256File("hashing.go")
	if err != nil {
		t.Error(err)
	}
}
