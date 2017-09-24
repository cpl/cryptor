package chunker

import (
	"bytes"
	"encoding/binary"

	"github.com/thee-engineer/cryptor/crypt"
)

// HeaderSize ...
const HeaderSize = 128

// ChunkHeader ...
type ChunkHeader struct {
	NKey crypt.AESKey // Key for the next chunk
	Hash []byte       // Hash of the chunk content
	Next []byte       // Hash of the next chunk
	Padd uint32       // Byte size of the padding
}

// NewChunkHeader ...
func NewChunkHeader() (header *ChunkHeader) {
	return &ChunkHeader{
		NKey: crypt.NullKey,
		Padd: 0,
		Next: nil,
		Hash: nil,
	}
}

// Bytes ...
func (header *ChunkHeader) Bytes() []byte {
	buffer := new(bytes.Buffer)

	buffer.Write(header.NKey.Bytes()) // 32
	buffer.Write(header.Next)         // 32
	buffer.Write(header.Hash)         // 32

	// Convert uint32 to byte array
	uintConv := make([]byte, 4)
	binary.LittleEndian.PutUint32(uintConv, header.Padd)
	buffer.Write(uintConv) // 32

	return buffer.Bytes()
}
