package chunker

import (
	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// Chunker ...
type Chunker struct {
	size   uint
	chunks []*chunk
	db     cachedb.Database // ! DEBUG
}

// New ...
func New(chunkSize uint, db cachedb.Database) *Chunker {
	return &Chunker{
		size:   chunkSize,
		chunks: []*chunk{newChunk(chunkSize)},
		db:     db,
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
// ! DEBUG
func (c *Chunker) Pack(key aes.Key) error {
	var nkey aes.Key
	next := make([]byte, 32)

	defer crypt.ZeroBytes(nkey[:], next)
	b := c.db.NewBatch()

	// Iterate chunks
	for _, chk := range c.chunks {
		nkey = aes.NewKey()

		// Pack the header and encrypt the chunk
		data, err := chk.pack(key, nkey, next[:])
		if err != nil {
			return err
		}
		b.Put(data)

		// Compute encrypted chunk hash for next hash
		next = hashing.Hash(data)
		key = nkey
	}

	return b.Write()
}
