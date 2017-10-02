// Package chunker contains data structures and functions that
// read, chunk and write from an io.Reader to a Database cache where
// data is stored in chunks (the hash acts as key).
package chunker

import (
	"io"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
)

// Chunker takes a reader as input, a chunk size and the cache in which to
// store any resulting chunks.
type Chunker struct {
	Size   uint32
	Cache  cachedb.Database
	Reader io.Reader
}

// Chunk starts chunking all data from the Chunker reader into the cache.
// If a non-null AES Key is given, the the tail chunk will be encrypted
// using this key, allowing the user more control. If a null key is given
// then a random AES Key will be used
func (c Chunker) Chunk(tKey crypt.AESKey) (pHash []byte, err error) {
	// Make a chunk struct
	chunk := NewChunk(c.Size)

	// Prepare previous hash and key
	pKey := crypt.NullKey    // Previous key is empty (this is the first chunk)
	pHash = make([]byte, 32) // Previous hash is empty

	// Prepare a batch for the cache, all chunks will be written at once
	batch := c.Cache.NewBatch()

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
		if read < int(c.Size) {
			chunk.Content = append(
				chunk.Content[:read],
				crypt.RandomData(uint(c.Size)-uint(read))...)
			chunk.Header.Padd = c.Size - uint32(read)
		} else {
			// No padding needed
			chunk.Header.Padd = 0
		}

		// Compute content hash for validity check
		chunk.Header.Hash = crypt.SHA256Data(chunk.Content[:read]).Sum(nil)

		// Store previous encryption key inside this chunk's header
		chunk.Header.NKey = pKey

		// Store previous encrypted chunk hash inside this chunk's header
		chunk.Header.Next = pHash

		// Generatea a new encryption key for each chunk
		// TODO: Find a better way of checking for last chunk
		if read < int(c.Size) {
			// Use tail key for the last chunk
			pKey = tKey
		} else {
			// Generate new random key
			pKey = crypt.NewKey()
		}

		// Encrypt chunk data
		eData, err := crypt.Encrypt(pKey, chunk.Bytes())
		if err != nil {
			return nil, err
		}

		// Hash encrypted content
		eHash := crypt.SHA256Data(eData).Sum(nil)

		// Store chunk in cache batch
		if err := batch.Put(eHash, eData); err != nil {
			return nil, err
		}

		// Update previous hash
		pHash = eHash
	}

	// Write batch to cache
	if err := batch.Write(); err != nil {
		return nil, err
	}

	// Return the tail hash
	return pHash, nil
}
