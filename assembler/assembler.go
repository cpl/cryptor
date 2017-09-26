package assembler

import (
	"encoding/binary"
	"errors"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

func extractHeader(data []byte) *chunker.ChunkHeader {
	// Check that given data is valid
	if len(data) < chunker.HeaderSize {
		panic(errors.New("Given chunk is too small"))
	}

	// Convert byte array to uint32
	padd := binary.LittleEndian.Uint32(data[:32])

	// Get next key from chunk header data
	nKey := crypt.NewKeyFromBytes(data[:32])

	return &chunker.ChunkHeader{
		NKey: nKey,
		Next: data[32:64],
		Hash: data[64:96],
		Padd: padd,
	}
}
