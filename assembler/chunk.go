package assembler

import (
	"encoding/binary"
	"errors"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

// EChunk ...
type EChunk struct {
	Data []byte
}

// GetEChunk returns the chunk with matching hash key from the given cache.
// TODO: Replace GetEChunk with Get call in Database. Move this as a method
// on Assembler struct.
// func (a *Assembler) GetChunk(hash) {
//   a.cache.Get(hash)
// }
func GetEChunk(hash []byte, cache cachedb.Database) *EChunk {

	// Get data from cache
	data, err := cache.Get(hash)
	if err != nil {
		panic(err)
	}

	return &EChunk{
		Data: data,
	}
}

// Decrypt returns the decrypted content of the encrypted chunk as a
// normal chunk containign a chunk header and content ([]byte).
func (eC EChunk) Decrypt(key crypt.AESKey) (*chunker.Chunk, error) {
	// Decrypt encrypted chunk data
	data, err := crypt.Decrypt(key, eC.Data)
	if err != nil {
		return nil, err
	}

	// Extract header from decrypted data
	header := extractHeader(data)

	// Return the chunk
	return &chunker.Chunk{
		Header:  header,
		Content: data[chunker.HeaderSize : len(data)-int(header.Padd)],
	}, nil
}

func extractHeader(data []byte) *chunker.ChunkHeader {
	// Check that given data is valid
	if len(data) < chunker.HeaderSize {
		panic(errors.New("Given chunk is too small"))
	}

	// Convert byte array to uint32
	padd := binary.LittleEndian.Uint32(data[96:100])

	// Get next key from chunk header data
	nKey := crypt.NewKeyFromBytes(data[:32])

	return &chunker.ChunkHeader{
		NKey: nKey,
		Next: data[32:64],
		Hash: data[64:96],
		Padd: padd,
	}
}
