package chunker

import (
	"bytes"
	"encoding/binary"

	"github.com/thee-engineer/cryptor/crypt"
)

// HeaderSize ...
const HeaderSize = 100

// ChunkHeader ...
type ChunkHeader struct {
	NKey crypt.AESKey // Key for the next chunk
	Next []byte       // Hash of the next chunk
	Hash []byte       // Hash of the chunk content
	Padd uint32       // Byte size of the padding
}

// NewChunkHeader ...
func NewChunkHeader() *ChunkHeader {
	return &ChunkHeader{
		NKey: crypt.NullKey,
		Next: nil,
		Hash: nil,
		Padd: 0,
	}
}

// Bytes ...
func (header ChunkHeader) Bytes() []byte {
	buffer := new(bytes.Buffer)

	buffer.Write(header.NKey.Bytes()) // 32
	buffer.Write(header.Next)         // 32
	buffer.Write(header.Hash)         // 32

	// Convert uint32 to byte array
	uintConv := make([]byte, 4)
	binary.LittleEndian.PutUint32(uintConv, header.Padd)
	buffer.Write(uintConv) // 4

	return buffer.Bytes()
}
