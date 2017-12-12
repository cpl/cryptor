package assembler

import "github.com/thee-engineer/cryptor/crypt/aes"

// Assembler ...
type Assembler interface {
	Assemble(key aes.Key, destination string) error
	// GetChunk(hash []byte) (EChunk, error)
}
