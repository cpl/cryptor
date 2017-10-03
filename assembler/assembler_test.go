package assembler

import (
	"bytes"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

func createTestCache() ([]byte, error) {
	// Create cache
	cache, err := cachedb.NewLDBCache("data", 0, 0)
	if err != nil {
		return nil, err
	}
	defer cache.Close()

	// Create archive and data buffers
	var buffer bytes.Buffer
	archive.TarGz("chunk.go", &buffer)

	// Create chunker
	c := &chunker.Chunker{
		Size:   1024,
		Cache:  cache,
		Reader: &buffer,
	}

	// Chunk files and get tail hash
	tail, err := c.Chunk(crypt.NullKey)
	if err != nil {
		return nil, err
	}

	return tail, nil
}

func TestAssembler(t *testing.T) {
	t.Parallel()

	// Create test cache
	tail, err := createTestCache()
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll("data")

	// Open cache
	cache, err := cachedb.NewLDBCache("data", 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer cache.Close()

	// Create assembler
	asm := &Assembler{
		Tail:  tail,
		Cache: cache,
	}

	// Start assembling package
	defer os.RemoveAll("untar")
	err = asm.Assemble(crypt.NullKey)
	if err != nil {
		t.Error(err)
	}
}
