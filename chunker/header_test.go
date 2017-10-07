package chunker_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

func TestChunkHeader(t *testing.T) {
	t.Parallel()

	// Create chunk header
	header := chunker.NewChunkHeader()

	header.Hash = crypt.RandomData(32) // Random content hash
	header.Next = crypt.RandomData(32) // Random next hash

	// Check header size to match expectations
	if len(header.Bytes()) != chunker.HeaderSize {
		t.Errorf("header error: invalid size; expected %d; got %d;",
			chunker.HeaderSize, len(header.Bytes()))
	}
}
