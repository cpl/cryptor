package assembler

import (
	"encoding/binary"
	"errors"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

// EChunk ...
type EChunk []byte

// GetEChunk returns the chunk with matching hash key from the given cache.
func GetEChunk(hash []byte, cache cachedb.Database) EChunk {

	// Get data from cache
	eChunk, err := cache.Get(hash)
	if err != nil {
		panic(err)
	}

	return eChunk
}

// Decrypt returns the decrypted content of the encrypted chunk as a
// normal chunk containign a chunk header and content ([]byte).
func (eChunk EChunk) Decrypt(key aes.Key) (*chunker.Chunk, error) {
	// Decrypt encrypted chunk data
	data, err := aes.Decrypt(key, eChunk)
	if err != nil {
		return nil, err
	}

	// Extract header from decrypted data
	header, err := extractHeader(data)
	if err != nil {
		return nil, err
	}

	// Return the chunk
	return &chunker.Chunk{
		Header:  header,
		Content: data[chunker.HeaderSize : len(data)-int(header.Padd)],
	}, nil
}

func extractHeader(data []byte) (*chunker.ChunkHeader, error) {
	// Check that given data is valid
	if len(data) < chunker.HeaderSize {
		return nil, errors.New("chunk extract header err | chunk is too small")
	}

	// Convert byte array to uint32
	padd := binary.LittleEndian.Uint32(data[96:100])

	// Get next key from chunk header data
	nKey, _ := aes.NewKeyFromBytes(data[:32])

	return &chunker.ChunkHeader{
		NKey: nKey,
		Next: data[32:64],
		Hash: data[64:96],
		Padd: padd,
	}, nil
}
