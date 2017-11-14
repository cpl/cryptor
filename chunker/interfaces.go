package chunker

import "github.com/thee-engineer/cryptor/crypt/aes"

// Chunker ...
type Chunker interface {
	Chunk(aes.Key) ([]byte, error)
}
