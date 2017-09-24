package chunker

import (
	"bytes"
	"encoding/binary"

	"github.com/thee-engineer/cryptor/crypt"
)

// HeaderSize ...
const HeaderSize = 96

// ChunkHeader ...
type ChunkHeader struct {
	Padd uint32
	NKey crypt.AESKey
	Next []byte
}

// NewChunkHeader ...
func NewChunkHeader() (header *ChunkHeader) {
	return &ChunkHeader{
		NKey: crypt.NullKey,
		Padd: 0,
		Next: nil,
	}
}

// Bytes ...
func (header *ChunkHeader) Bytes() []byte {
	buffer := new(bytes.Buffer)

	buffer.Write(header.Next)         // 32
	buffer.Write(header.NKey.Bytes()) // 32

	uintConv := make([]byte, 4)
	binary.LittleEndian.PutUint32(uintConv, header.Padd)
	buffer.Write(uintConv) // 32

	return buffer.Bytes()
}
