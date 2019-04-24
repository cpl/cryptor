package mwords

import (
	"hash/crc32"
	"testing"
)

// the total number of characters in the word list
const charCount = 11068

// expected checksum
const checksum = 2176441764

var wordsValid = []string{
	"apple", "approve", "canvas", "cruise", "fame",
	"merry", "salad", "soda", "trust", "zoo",
}

var wordsInvalid = []string{
	"applezoo", "", "", "søda", "sodaa", "sod", "01234", "salad_", ".salad",
}

func TestIsValidWord(t *testing.T) {
	t.Parallel()

	// test a set of valid words
	for _, word := range wordsValid {
		if !IsValidWord(word) {
			t.Errorf("expected %s to be valid word\n", word)
		}
	}

	// test a set of invalid words
	for _, word := range wordsInvalid {
		if IsValidWord(word) {
			t.Errorf("expected %s to be invalid word\n", word)
		}
	}
}

func TestValidateWordlist(t *testing.T) {
	t.Parallel()

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

		// at the same time test the lookup map
		if idx, ok := mnemonicLookup[word]; !ok {
			t.Fatalf("failed to find expected word %s in lookup map\n", word)
		} else {
			if mnemonicWords[idx] != word {
				t.Fatalf("mismatch word at index %d, got %s but want %s\n",
					idx, mnemonicWords[idx], word)
			}
		}
	}

	// perform CRC32 checksum on data
	if crc32.ChecksumIEEE(totalData) != checksum {
		t.Fatal("invalid checksum")
	}
}
