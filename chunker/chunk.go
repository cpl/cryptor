package chunker

import (
	"bytes"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// NullByteArray is used for the last chunk header.Next
var NullByteArray [32]byte

// Chunk combines the chunk content which is a []byte of header.Size and the
// chunk header which contains information about the this chunk and next chunk.
type Chunk struct {
	Header  *ChunkHeader
	Content []byte
}

// NewChunk creates a new chunk with given size content
func NewChunk(size uint32) *Chunk {
	return &Chunk{
		Header:  NewChunkHeader(),
		Content: make([]byte, size),
	}
}

// Bytes returns the chunk header and content as []byte
func (c Chunk) Bytes() []byte {
	var buffer bytes.Buffer

	buffer.Write(c.Header.Bytes()) // Write header bytes
	buffer.Write(c.Content)        // Write chunk content

	return buffer.Bytes()
}

// IsValid compares the header hash with the content hash
func (c Chunk) IsValid() bool {
	return bytes.Equal(c.Header.Hash, hashing.SHA256Digest(c.Content))
}

// IsLast checks if the next chunk hash is the NullByteArray
func (c Chunk) IsLast() bool {
	return bytes.Equal(c.Header.Next, NullByteArray[:])
}

// Zero writes all the chunk content with 0 and the chunk header data
func (c Chunk) Zero() {
	crypt.ZeroBytes(c.Content)
	c.Header.Zero()
}

func (c *Chunk) setHeader(nextKey aes.Key, nextHash []byte, read int) {
	// Compute content hash for validity check
	c.Header.Hash = hashing.SHA256Digest(c.Content[:read])

	// Store previous encryption key inside this chunk's header
	c.Header.NKey = nextKey

	// Store previous encrypted chunk hash inside this chunk's header
	c.Header.Next = nextHash
}

func (c *Chunk) padd(realSize int) {
	expectedSize := cap(c.Content)

	// Add random padding if needed
	if realSize < expectedSize {
		c.Content = append(
			c.Content[:realSize],
			crypt.RandomData(uint(expectedSize)-uint(realSize))...)
		c.Header.Padd = uint32(expectedSize) - uint32(realSize)
	} else {
		// No padding needed
		c.Header.Padd = 0
	}
}
