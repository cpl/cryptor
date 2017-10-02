package chunker

import "testing"
import "github.com/thee-engineer/cryptor/crypt"

func TestChunkHeader(t *testing.T) {
	t.Parallel()

	// Create chunk header
	header := NewChunkHeader()

	header.Hash = crypt.RandomData(32) // Random content hash
	header.Next = crypt.RandomData(32) // Random next hash

	// Check header size to match expectations
	if len(header.Bytes()) != HeaderSize {
		t.Errorf("header error: invalid size; expected %d; got %d;",
			HeaderSize, len(header.Bytes()))
	}
}
