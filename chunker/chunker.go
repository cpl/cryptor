// Package chunker contains data structures and functions that
// read, chunk and write from an io.Reader to a Database cache where
// data is stored in chunks (the hash acts as key).
package chunker

import (
	"io"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// DefaultChunker takes a reader as input, a chunk size and the cache in which to
// store any resulting chunks.
type DefaultChunker struct {
	ChunkSize uint32
	Cache     cachedb.Database
	Reader    io.Reader
}

// NewDefaultChunker ...
func NewDefaultChunker(Reader io.Reader, ChunkSize uint32, Cache cachedb.Database) Chunker {
	return &DefaultChunker{ChunkSize: ChunkSize, Cache: Cache, Reader: Reader}
}

// Chunk starts chunking all data from the Chunker reader into the cache.
// If a non-null AES Key is given, the the tail chunk will be encrypted
// using this key, allowing the user more control. If a null key is given
// then a random AES Key will be used.
// TODO: Rnadomize chunk order (read all data first, count it, shuffle it,
// padd data if needed, encrypt each chunk as normal)
func (c DefaultChunker) Chunk(tailKey aes.Key) (tailHash []byte, err error) {
	// Make a chunk struct
	chunk := NewChunk(c.ChunkSize)

	// Prepare previous hash and key
	pKey := aes.NullKey         // Previous key is empty (first chunk)
	tailHash = make([]byte, 32) // Previous hash is empty (at the end, tail)

	// Prepare a batch for the cache, all chunks will be written at once
	batch := c.Cache.NewBatch()

	// Zero memory of tail key, previous key and tail hash after chunking
	defer crypt.ZeroBytes(tailKey[:])
	defer crypt.ZeroBytes(pKey[:])
	defer crypt.ZeroBytes(tailHash[:])

	for {
		// Read archive content into chunks
		read, err := c.Reader.Read(chunk.Content)

		// Check for EOF
		if read == 0 || err == io.EOF {
			break
		}

		// Check for errors
		if err != nil {
			return nil, err
		}

		// Add random padding if needed
		if read < int(c.ChunkSize) {
			chunk.Content = append(
				chunk.Content[:read],
				crypt.RandomData(uint(c.ChunkSize)-uint(read))...)
			chunk.Header.Padd = c.ChunkSize - uint32(read)
		} else {
			// No padding needed
			chunk.Header.Padd = 0
		}

		// Compute content hash for validity check
		chunk.Header.Hash = hashing.SHA256Digest(chunk.Content[:read])

		// Store previous encryption key inside this chunk's header
		chunk.Header.NKey = pKey

		// Store previous encrypted chunk hash inside this chunk's header
		chunk.Header.Next = tailHash

		// Generatea a new encryption key for each chunk
		// TODO: Find a better way of checking for last chunk
		if read < int(c.ChunkSize) {
			// Use tail key for the last chunk
			pKey = tailKey
		} else {
			// Generate new random key
			pKey = aes.NewKey()
		}

		// Encrypt chunk data
		eData, err := aes.Encrypt(pKey, chunk.Bytes())
		if err != nil {
			return nil, err
		}

		// Zero key from memory
		crypt.ZeroBytes(chunk.Header.NKey[:])

		// Hash encrypted content
		eHash := hashing.SHA256Digest(eData)

		// Store chunk in cache batch
		if err := batch.Put(eHash, eData); err != nil {
			return nil, err
		}

		// Update previous hash
		tailHash = eHash
	}

	// Write batch to cache
	if err := batch.Write(); err != nil {
		return nil, err
	}

	// Return the tail hash
	return tailHash, nil
}
