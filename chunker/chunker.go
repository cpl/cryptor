package chunker

import (
	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
	"github.com/thee-engineer/cryptor/crypt/scrypt"
)

// Chunker ...
type Chunker struct {
	size    uint
	chunks  []*Chunk
	manager cachedb.Manager
}

// New ...
func New(chunkSize uint, manager cachedb.Manager) *Chunker {
	return &Chunker{
		size:    chunkSize,
		chunks:  []*Chunk{newChunk(chunkSize)},
		manager: manager,
	}
}

// Write ...
func (c *Chunker) Write(p []byte) (n int, err error) {
	// Write data to current chunk
	n, nerr := c.chunks[len(c.chunks)-1].Write(p)
	// Append more chunks if data does not fit
	if nerr != nil && nerr.Error() == "data does not fit inside chunk" {
		c.chunks = append(c.chunks, newChunk(c.size))
	}
	return n, nil
}

// Pack ...
func (c *Chunker) Pack(password string) (tail []byte, err error) {
	var key, nkey aes.Key
	nhash := make([]byte, hashing.HashSize)

	// Erase keys
	defer crypt.ZeroBytes(key[:], nkey[:])

	// Iterate chunks
	for index, chk := range c.chunks {

		// If "tail" chunk, encrypt using password
		if index == len(c.chunks)-1 {
			key, err = aes.NewKeyFromBytes(scrypt.Scrypt(password, []byte{}))
			if err != nil {
				return nil, nil
			}
		} else {
			// Use random key for any other chunk
			key = aes.NewKey()
		}

		// Encrypt chunk with key and append previous chunk key and hash
		data, err := chk.pack(key, nkey, nhash)
		if err != nil {
			return nil, err
		}

		// Add encrypted chunk to cache
		if err := c.manager.Add(data); err != nil {
			return nil, err
		}

		// Store previous chunk hash and key
		nhash = hashing.Hash(data)
		nkey = key
	}

	return nhash, nil
}
