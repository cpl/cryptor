package chunker

import "testing"
import "github.com/thee-engineer/cryptor/crypt"

func TestChunk(t *testing.T) {
	t.Parallel()

	// Create a chunk of 1024 bytes
	chunk := NewChunk(1024)
	// Fill the chunk with random data
	chunk.Content = crypt.RandomData(1024)

	// Create the chunk header
	chunkHeader := NewChunkHeader()
	chunkHeader.Hash = crypt.SHA256Data(chunk.Content).Sum(nil) // Content hash
	chunkHeader.Next = NullByteArray[:]                         // Tail hash
	chunkHeader.NKey = crypt.NullKey                            // Tail key

	chunk.Header = chunkHeader

	if !chunk.IsValid() {
		t.Error("chunk error: invalid chunk")
	}

	if !chunk.IsLast() {
		t.Error("chunk error: expected last")
	}

	chunk.Bytes()
}
