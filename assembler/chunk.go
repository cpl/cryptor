package assembler

import (
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

// EChunk ...
type EChunk struct {
	Data []byte
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
