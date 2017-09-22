package chunker

import (
	"bytes"
	"encoding/binary"
)

// HeaderSize ...
const HeaderSize = 128

// ChunkHeader ...
type ChunkHeader struct {
	Hash []byte
	Padd uint32
	NKey []byte
	Prev []byte
}

// NewChunkHeader ...
func NewChunkHeader() (header *ChunkHeader) {
	return &ChunkHeader{
		Hash: nil,
		NKey: nil,
		Padd: 0,
		Prev: nil,
	}
}

// Bytes ...
func (header *ChunkHeader) Bytes() []byte {
	buffer := new(bytes.Buffer)

	buffer.Write(header.Hash) // 32
	buffer.Write(header.Prev) // 32
	buffer.Write(header.NKey) // 32

	uintConv := make([]byte, 4)
	binary.LittleEndian.PutUint32(uintConv, header.Padd)
	buffer.Write(uintConv) // 32

	return buffer.Bytes()
}
