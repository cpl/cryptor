package chunker

import (
	"bytes"
	"encoding/binary"

	"github.com/thee-engineer/cryptor/crypt"

	"github.com/thee-engineer/cryptor/crypt/aes"
)

// HeaderSize defines the expected size of the header
const HeaderSize = 100

// ChunkHeader ...
type ChunkHeader struct {
	NKey aes.Key // Key for the next chunk
	Next []byte  // Hash of the next chunk
	Hash []byte  // Hash of the chunk content
	Padd uint32  // Byte size of the padding
}

// NewChunkHeader returns a new chunk header
func NewChunkHeader() *ChunkHeader {
	return &ChunkHeader{
		NKey: aes.NullKey,
		Next: nil,
		Hash: nil,
		Padd: 0,
	}
}

// Bytes returns the chunk header as a byte array. Data is distributed as
// follows:
//
// | NKEY 32B | NEXT 32B | HASH 32B | PADD 4B |
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

// Zero writes all header data with 0
func (header ChunkHeader) Zero() {
	crypt.ZeroBytes(header.Hash[:], header.Next[:], header.NKey[:])
	header.Padd = 0
}
