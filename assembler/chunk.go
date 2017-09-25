package assembler

import (
	"bytes"
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

func getEChunk(hash string) (*EChunk, error) {
	// Obtain chunk path
	cachePath := cache.GetCachePath()
	chunkPath := path.Join(cachePath, hash)

	// Check that chunks exist
	_, err := os.Stat(chunkPath)
	if err != nil {
		return nil, err
	}

	// Obtain chunk contents
	data, err := ioutil.ReadFile(chunkPath)
	if err != nil {
		return nil, err
	}

	return &EChunk{
		Data: data,
	}, nil
}

// Chunk ...
type Chunk struct {
	Header  *chunker.ChunkHeader
	Content []byte
}

// GetChunk ...
func GetChunk(hash string, key *crypt.AESKey) (*Chunk, error) {
	eChunk, err := getEChunk(hash)
	if err != nil {
		return nil, err
	}

	return eChunk.Decrypt(key), nil
}

// Decrypt ...
func (eC *EChunk) Decrypt(key *crypt.AESKey) *Chunk {
	// Decrypt eChunk
	data, err := crypt.Decrypt(key, eC.Data)
	if err != nil {
		panic(err)
	}

	// Extract header
	header := extractHeader(data)
	// Extract content
	content := data[chunker.HeaderSize : len(data)-int(header.Padd)]

	// Validate chunk content
	if bytes.Compare(header.Hash, crypt.SHA256Data(content).Sum(nil)) != 0 {
		panic(errors.New("Mismatch header hash with content hash"))
	}

	return &Chunk{
		Header:  header,
		Content: content,
	}
}

func extractHeader(data []byte) *chunker.ChunkHeader {
	// Check that given data is valid
	if len(data) < chunker.HeaderSize {
		panic(errors.New("Given chunk is too small"))
	}

	// Convert byte array to uint32
	padd := binary.LittleEndian.Uint32(data[:32])

	// Get next key from chunk header
	var nKey crypt.AESKey
	copy(nKey[:], data[:32])

	return &chunker.ChunkHeader{
		NKey: nKey,
		Next: data[32:64],
		Hash: data[64:96],
		Padd: padd,
	}
}
