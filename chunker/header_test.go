package chunker

import "testing"
import "github.com/thee-engineer/cryptor/crypt"

func TestChunkHeader(t *testing.T) {
	header := NewChunkHeader()

	header.Hash = crypt.RandomData(32)
	header.Next = crypt.RandomData(32)

	if len(header.Bytes()) != HeaderSize {
		t.Errorf("Invalid header size! Expected %d, got %d",
			HeaderSize, len(header.Bytes()))
	}
}
