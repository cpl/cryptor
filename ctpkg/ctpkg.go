package ctpkg

import (
	"bytes"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

// CTPKG ...
type CTPKG struct {
	Name        string
	Hash        string
	Size        int
	ChunkCount  int
	ChunkHashes []string
	PKey        string
}

// NewCTPKG ...
func NewCTPKG(source, name string, chunkSize uint32) *CTPKG {
	contentBuffer := new(bytes.Buffer)

	// Create tar.gz of file/dir
	if err := archive.TarGz(source, contentBuffer); err != nil {
		panic(err)
	}

	// Hash tar.gz for integrity check
	contentHash := crypt.SHA256Data(contentBuffer.Bytes())

	// Generate a random primary key for the package
	pKey := crypt.NewKey()

	// Create a chunker
	chunker := &chunker.Chunker{
		Size:   chunkSize,
		Reader: contentBuffer,
		Key:    pKey,
	}

	// Start chunking the tar.gz
	hashList, err := chunker.Chunk()
	if err != nil {
		panic(err)
	}

	// Create package info
	ctpkg := &CTPKG{
		Name:        name,
		Hash:        string(crypt.Encode(contentHash.Sum(nil))),
		Size:        contentBuffer.Len(),
		PKey:        string(crypt.Encode(pKey[:])),
		ChunkCount:  len(hashList),
		ChunkHashes: hashList,
	}

	return ctpkg
}

// ToJSON ...
func (ctpkg *CTPKG) ToJSON() error {
	return nil
}
