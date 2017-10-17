package chunker_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

func TestChunk(t *testing.T) {
	t.Parallel()

	// Create a chunk of 1024 bytes
	chunk := chunker.NewChunk(1024)
	// Fill the chunk with random data
	chunk.Content = crypt.RandomData(1024)

	// Create the chunk header
	chunkHeader := chunker.NewChunkHeader()
	chunkHeader.Hash = hashing.SHA256Digest(chunk.Content) // Content hash
	chunkHeader.Next = chunker.NullByteArray[:]            // Tail hash
	chunkHeader.NKey = aes.NullKey                         // Tail key

	chunk.Header = chunkHeader

	if !chunk.IsValid() {
		t.Error("chunk error: invalid chunk")
	}

	if !chunk.IsLast() {
		t.Error("chunk error: expected last")
	}

	chunk.Bytes()
}
