package chunk

import (
	"bytes"
	"errors"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// NullByteArray is used for the last chunk header.Next
var NullByteArray [32]byte

// Chunk ...
type Chunk struct {
	Head *header
	Body []byte

	size int
}

// New ...
func New(size uint) *Chunk {
	return &Chunk{
		Head: newHeader(),
		Body: make([]byte, size),
		size: 0,
	}
}

// Unpack ...
// TODO Finish for extraction/decryption/assembly
func Unpack(key aes.Key, data []byte) *Chunk {
	return nil
}

// Bytes ...
func (c *Chunk) Bytes() []byte {
	data := make([]byte, HeaderSize+cap(c.Body))
	copy(data[0:HeaderSize], c.Head.Bytes())
	copy(data[HeaderSize:], c.Body)
	return data
}

// Zero ...
func (c *Chunk) Zero() {
	crypt.ZeroBytes(c.Body)
	c.Head.Zero()
}

// Read ...
// TODO Finish for extraction/decryption/assembly
func (c *Chunk) Read(p []byte) (n int, err error) {
	return 0, nil
}

// IsValid ...
// TODO Finish for extraction/decryption/assembly
func (c *Chunk) IsValid() bool {
	return bytes.Equal(c.Head.Hash, hashing.Hash(c.Body[:c.size]))
}

// IsLast ...
// TODO Finish for extraction/decryption/assembly
func (c *Chunk) IsLast() bool {
	return bytes.Equal(c.Head.NextHash, NullByteArray[:])
}

// Write ...
func (c *Chunk) Write(p []byte) (n int, err error) {
	// Check if write exceeded chunk body size
	if c.size+len(p) > cap(c.Body) {
		return 0, errors.New("data does not fit inside chunk")
	}

	// Copy data inside chunk
	copy(c.Body[c.size:c.size+len(p)], p)
	c.size += len(p)

	// Recompute hash of content
	// c.Head.Hash = hashing.Hash(c.Body[:c.size])

	// Update padding
	// if c.Head.Padding > 0 {
	// 	c.Head.Padding -= uint32(len(p))
	// }

	return len(p), nil
}

// Padd ...
func (c *Chunk) padd() {
	// Check if chunk is full
	if c.size == cap(c.Body) {
		c.Head.Padding = 0
		return
	}

	// Calculate required padding
	c.Head.Padding = uint32(cap(c.Body)) - uint32(c.size)

	// Add random padding to chunk
	copy(c.Body[c.size:], crypt.RandomData(uint(c.Head.Padding)))
}

// Pack ...
func (c *Chunk) Pack(key, nkey aes.Key, next []byte) ([]byte, error) {
	// Remove chunk after packing
	defer c.Zero()
	defer crypt.ZeroBytes(key[:], nkey[:], next)

	// Compute content hash
	c.Head.Hash = hashing.Hash(c.Body[:c.size])
	// Add padding
	c.padd()
	// Add next chunk hash & key
	c.Head.NextHash = next
	c.Head.NextKey = nkey

	// Encrypt chunk
	e, err := aes.Encrypt(key, c.Bytes())
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Decrypt ...
func Decrypt(key aes.Key, data []byte) (*Chunk, error) {
	return nil, nil
}
