package chunker

import (
	"bytes"
	"errors"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// NullByteArray is used for the last chunk header.Next
var NullByteArray [32]byte

type chunk struct {
	Head *header
	Body []byte

	size int
}

func newChunk(size uint) *chunk {
	return &chunk{
		Head: newHeader(),
		Body: make([]byte, size),
		size: 0,
	}
}

// Bytes ...
func (c *chunk) Bytes() []byte {
	data := make([]byte, HeaderSize+cap(c.Body))
	copy(data[0:HeaderSize], c.Head.Bytes())
	copy(data[HeaderSize:], c.Body)
	return data
}

// Zero ...
func (c *chunk) Zero() {
	crypt.ZeroBytes(c.Body)
	c.Head.Zero()
}

// Read ...
func (c *chunk) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (c *chunk) isValid() bool {
	return bytes.Equal(c.Head.Hash, hashing.Hash(c.Body[:c.size]))
}

func (c *chunk) isLast() bool {
	return bytes.Equal(c.Head.NextHash, NullByteArray[:])
}

// Write ...
func (c *chunk) Write(p []byte) (n int, err error) {
	// Check if write exceeded chunk body size
	if c.size+len(p) > cap(c.Body) {
		return 0, errors.New("data does not fit inside chunk")
	}

	// Copy data inside chunk
	copy(c.Body[c.size:c.size+len(p)], p)
	c.size += len(p)

	return len(p), nil
}

// Padd ...
func (c *chunk) padd() {
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

func (c *chunk) pack(key, nkey aes.Key, nhash []byte) ([]byte, error) {
	// Remove chunk after packing
	defer c.Zero()
	defer crypt.ZeroBytes(key[:], nkey[:], nhash)

	// Compute content hash
	c.Head.Hash = hashing.Hash(c.Body[:c.size])
	// Add padding
	c.padd()
	// Add next chunk hash & key
	c.Head.NextHash = nhash
	c.Head.NextKey = nkey

	// Encrypt chunk
	e, err := aes.Encrypt(key, c.Bytes())
	if err != nil {
		return nil, err
	}
	return e, nil
}
