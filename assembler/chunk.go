package assembler

import (
	"encoding/binary"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

func loadChunkHeader(data []byte) *chunker.ChunkHeader {
	// Get data from chunk
	// data := c.Data

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
