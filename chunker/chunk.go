package chunker

import (
	"bytes"

	"github.com/thee-engineer/cryptor/crypt"
)

// NullByteArray ...
var NullByteArray [32]byte

// Chunk ...
type Chunk struct {
	Header  *ChunkHeader
	Content []byte
}

// NewChunk ...
func NewChunk(size uint32) *Chunk {
	return &Chunk{
		Header:  NewChunkHeader(),
		Content: make([]byte, size),
	}
}

// Bytes ...
func (c Chunk) Bytes() []byte {
	var buffer bytes.Buffer

	buffer.Write(c.Header.Bytes())
	buffer.Write(c.Content)

	return buffer.Bytes()
}

// IsValid ...
func (c Chunk) IsValid() bool {
	return bytes.Equal(c.Header.Hash, crypt.SHA256Data(c.Content).Sum(nil))
}

// IsLast ...
func (c Chunk) IsLast() bool {
	return bytes.Equal(c.Header.Next, NullByteArray[:])
}
