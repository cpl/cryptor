package chunker

import "bytes"
import "github.com/thee-engineer/cryptor/crypt"

// NullByteArray ...
var NullByteArray []byte

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
	return bytes.Compare(c.Header.Hash, crypt.SHA256Data(c.Content).Sum(nil)) == 0
}

// IsLast ...
func (c Chunk) IsLast() bool {
	return bytes.Compare(c.Header.Next, NullByteArray) == 0
}
