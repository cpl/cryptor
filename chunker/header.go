package chunker

// ChunkHeader ...
type ChunkHeader struct {
	Hash    []byte
	Size    int
	PadSize int
	NKey    []byte
}

// NewChunkHeader ...
func NewChunkHeader() (header *ChunkHeader) {
	return header
}

// Bytes ...
func (header *ChunkHeader) Bytes() []byte {
	return append(header.NKey[:], header.Hash[:]...)
}
