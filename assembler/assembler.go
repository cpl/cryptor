package assembler

import (
	"errors"
	"io/ioutil"
	"path"

	"github.com/thee-engineer/cryptor/cache"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

// Assembler ...
type Assembler struct {
	Size uint32
}

// Assemble ...
func (assembler *Assembler) Assemble(tail string, pKey crypt.AESKey) error {
	return nil
}

func (assembler *Assembler) getChunk(hash string) error {
	// Get chunk cache abs path
	chunkCachePath := cache.GetCachePath()
	// Look for chunk file in cache
	chunkFilePath := path.Join(chunkCachePath, hash)

	// Open chunk file for reading
	chunkData, err := ioutil.ReadFile(chunkFilePath)
	if err != nil {
		return err
	}

	// Invalid chunk file
	if len(chunkData) != int(assembler.Size)+chunker.HeaderSize {
		return errors.New("Invalid chunk file size")
	}

	return nil
}
