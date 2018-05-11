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

// NewDefaultChunker ...
// func NewDefaultChunker(Reader io.Reader, ChunkSize uint32, Cache cachedb.Database) Chunker {
// 	return &DefaultChunker{ChunkSize: ChunkSize, Cache: Cache, Reader: Reader}
// }

// ChunkFrom starts chunking all data from the Reader into the Cache.
// If a non-null AES Key is given, the the tail chunk will be encrypted
// using this key, allowing the user more control. If a null key is given
// then a random AES Key will be used.
func ChunkFrom(reader io.Reader, size uint32, cache cachedb.Manager, tailKey aes.Key) (nextHash []byte, err error) {
	// Make a chunk struct
	chunk := NewChunk(size)

	nextKey := aes.NullKey      // Next key is empty (first chunk)
	nextHash = make([]byte, 32) // Next hash is empty (at the end, tail)

	// Prepare a batch for the cache, all chunks will be written at once
	// batch := cache.NewBatch()

	// Zero memory of tail key, next key and tail hash after chunking
	defer crypt.ZeroBytes(tailKey[:], nextKey[:], nextHash[:])
	// Zero memory of the chunk struct
	defer chunk.Zero()

	for {
		// Read archive content into chunks
		read, err := reader.Read(chunk.Content)

		// Check for EOF
		if read == 0 || err == io.EOF {
			break
		}

		// Check for errors
		if err != nil {
			return nil, err
		}

		// Set the header (key, hash, next)
		chunk.setHeader(nextKey, nextHash, read)
		// Setup padding of chunk and update header
		chunk.padd(read)

		// Generatea a new encryption key for each chunk
		if read < int(size) {
			// Use tail key for the last chunk
			nextKey = tailKey
		} else {
			// Generate new random key
			nextKey = aes.NewKey()
		}

		// Encrypt chunk data
		encryptedData, err := aes.Encrypt(nextKey, chunk.Bytes())
		if err != nil {
			return nil, err
		}

		// Zero key from memory
		crypt.ZeroBytes(chunk.Header.NKey[:])

		// Hash encrypted content
		eHash := hashing.Hash(encryptedData)

		// Store chunk in cache batch
		if err := cache.Add(encryptedData); err != nil {
			return nil, err
		}

		// Update previous hash
		nextHash = eHash
	}

	// Write batch to cache
	// if err := batch.Write(); err != nil {
	// 	return nil, err
	// }

	// Return the tail hash
	return nextHash, nil
}

// func randomizePackageOrder(reader io.Reader) ([][]byte, error) {
// 	// Read all the package data
// 	packageContent, err := ioutil.ReadAll(reader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }
