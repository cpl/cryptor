package chunker

import "testing"
import "github.com/thee-engineer/cryptor/crypt"

func TestChunk(t *testing.T) {
	t.Parallel()

	chunk := NewChunk(1024)
	chunk.Content = crypt.RandomData(1024)

	chunkHeader := NewChunkHeader()
	chunkHeader.Hash = crypt.SHA256Data(chunk.Content).Sum(nil)
	chunkHeader.Next = NullByteArray[:]
	chunkHeader.NKey = crypt.NullKey

	chunk.Header = chunkHeader

	if !chunk.IsValid() {
		t.Error("Chunk is not valid!")
	}

	if !chunk.IsLast() {
		t.Error("Chunk is not last!")
	}

	chunk.Bytes()
}
