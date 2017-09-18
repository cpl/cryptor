package ctpkgm

// CTPKG ...
type CTPKG struct {
	Name        string
	Hash        string
	Size        uint
	ChunkCount  uint
	ChunkHashes []string
	Key         string
}

// NewCTPKG ...
func NewCTPKG(source, name string) (ctpkg CTPKG) {

	return ctpkg
}
