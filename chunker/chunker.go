package chunker

import (
	"github.com/thee-engineer/cryptor/cachedb"
)

// Chunker ...
type Chunker struct {
	size    uint
	chunks  []*chunk
	manager cachedb.Manager
}

// New ...
func New(chunkSize uint, manager cachedb.Manager) *Chunker {
	return &Chunker{
		size:    chunkSize,
		chunks:  []*chunk{newChunk(chunkSize)},
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
