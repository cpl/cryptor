package chunker

import (
	"bytes"
	"io"

	"github.com/thee-engineer/cryptor/crypt"
)

// Chunker ...
type Chunker struct {
	Size   uint32
	reader io.Reader
}

// NewChunker ...
func NewChunker(source io.Reader, size uint32) (chunker *Chunker) {
	chunker.Size = size
	chunker.reader = source

	return chunker
}

// Start ...
func (chunker *Chunker) Start() error {
	chunkHeader := NewChunkHeader()

	chunkData := new(bytes.Buffer)
	chunkCont := make([]byte, chunker.Size)

	for {
		// Read archive content into chunks
		read, err := chunker.reader.Read(chunkCont)

		// EOF
		if read == 0 || err != io.EOF {
			break
		}

		// Check for errors
		if err != nil {
			return err
		}

		// Update chunk header
		chunkHeader.Hash = crypt.SHA256Data(chunkCont).Sum(nil)
		chunkHeader.NKey = crypt.NewKey()[:]
		chunkHeader.Padd = chunker.Size - uint32(read)

		// Add padding if needed
		if read < int(chunker.Size) {
			for index := read; index < int(chunker.Size); index++ {
				chunkCont[index] = 0
			}
		}

		// Create chunk with header and content
		chunkData.Write(chunkHeader.Bytes())
		chunkData.Write(chunkCont)

	}

	return nil
}
