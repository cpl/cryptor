package mwords

import (
	"hash/crc32"
	"testing"
)

// the total number of characters in the word list
const charCount = 11068

// expected checksum
const checksum = 2176441764

func TestValidateWordlist(t *testing.T) {
	// validate expected length
	if len(mnemonicWords) != Count {
		t.Fatalf("invalid wordlist length, expected %d, got %d",
			Count, len(mnemonicWords))
	}

	// concat all words in a single byte array
	totalData := make([]byte, charCount)
	totalCount := 0
	for _, word := range mnemonicWords {
		copy(totalData[totalCount:], []byte(word))
		totalCount += len(word)
	}

	// perform CRC32 checksum on data
	if crc32.ChecksumIEEE(totalData) != checksum {
		t.Fatal("invalid checksum")
	}
}
