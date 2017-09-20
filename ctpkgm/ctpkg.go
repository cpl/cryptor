package ctpkgm

import (
	"bytes"
	"fmt"

	"github.com/thee-engineer/cryptor/chunker"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/crypto"
)

// CTPKG ...
type CTPKG struct {
	Name string
	Hash string
	Size int

	ChunkCount  int
	ChunkHashes []string

	Key string
}

// NewCTPKG ...
func NewCTPKG(source, name string) (ctpkg CTPKG) {
	contentBuffer := new(bytes.Buffer)

	if err := archive.TarGz(source, contentBuffer); err != nil {
		panic(err)
	}

	contentHash := crypto.SHA256Data(contentBuffer.Bytes())

	mainKey := crypto.NewKey()

	ctpkg.Name = name
	ctpkg.Hash = string(crypto.Encode(contentHash.Sum(nil)))
	ctpkg.Size = contentBuffer.Len()
	ctpkg.Key = string(crypto.Encode(mainKey[:]))

	chunkCount, chunksPath, err := chunker.ChunkData(contentBuffer.Bytes(), 1024, mainKey)
	if err != nil {
		panic(err)
	}

	ctpkg.ChunkCount = chunkCount

	return ctpkg
}

func (ctpkg *CTPKG) String() string {
	return fmt.Sprintf("Name: %s\nHash: %s\nSize: %d",
		ctpkg.Name, ctpkg.Hash, ctpkg.Size)
}
