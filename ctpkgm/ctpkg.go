package ctpkgm

import (
	"fmt"
	"os"
)

// CTPKG ...
type CTPKG struct {
	Name        string
	Hash        string
	Size        uint
	ChunkCount  uint
	ChunkHashes [][]byte
	Key         []byte
}

// NewCTPKG ...
func NewCTPKG(source string) CTPKG {
	fmt.Println(os.Stat(source))
	return CTPKG{"", "", 0, 0, nil, nil}
}
