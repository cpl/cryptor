package ctpkgm

import (
	"bytes"
	"fmt"

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

	contentHash, err := crypto.SHA256Data(contentBuffer.Bytes())
	if err != nil {
		panic(err)
	}

	ctpkg.Name = name
	ctpkg.Hash = string(crypto.Encode(contentHash.Sum(nil)))
	ctpkg.Size = contentBuffer.Len()
	ctpkg.Key = string(crypto.Encode(crypto.NewKey()[:]))

	return ctpkg
}

func (ctpkg *CTPKG) String() string {
	return fmt.Sprintf("Name: %s\nHash: %s\nSize: %d",
		ctpkg.Name, ctpkg.Hash, ctpkg.Size)
}
