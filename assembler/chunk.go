package assembler

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/cache"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

// EChunk ...
type EChunk struct {
	Data []byte
}

// GetEChunk ...
func GetEChunk(hash string) *EChunk {
	// Get chunk path and check that it exists
	eChunkPath := path.Join(cache.GetCachePath(), hash)
	if _, err := os.Stat(eChunkPath); os.IsNotExist(err) {
		panic(err)
	}

	// Get chunk content
	data, err := ioutil.ReadFile(eChunkPath)
	if err != nil {
		panic(err)
	}

	return &EChunk{
		Data: data,
	}
}

// Decrypt ...
func (eC EChunk) Decrypt(key crypt.AESKey) *chunker.Chunk {
	data, err := crypt.Decrypt(key, eC.Data)
	if err != nil {
		panic(err)
	}

	header := extractHeader(data)

	return &chunker.Chunk{
		Header:  header,
		Content: data[chunker.HeaderSize : len(data)-int(header.Padd)],
	}
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
