package assembler_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/assembler"
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
		Size:   16,
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
	asm := &assembler.Assembler{
		Tail:  tail,
		Cache: cache,
	}

	// Start assembling package
	defer os.RemoveAll("untar")
	err = asm.Assemble(crypt.NullKey, "untar")
	if err != nil {
		t.Error(err)
	}
}

func TestFullChunkAssemble(t *testing.T) {
	t.Parallel()

	// Prepare tar archive of cryptor in buffer
	var buffer bytes.Buffer
	if err := archive.TarGz("..", &buffer); err != nil {
		t.Error(err)
	}

	// Create cache for chunks
	cache, err := cachedb.NewLDBCache("/tmp/asmcnktest", 16, 16)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll("/tmp/asmcnktest")

	// Create chunker
	cnk := &chunker.Chunker{
		Size:   1000000,
		Cache:  cache,
		Reader: &buffer,
	}

	// Chunk with derived key
	key := crypt.NewKeyFromPassword("testing")
	tail, err := cnk.Chunk(key)
	if err != nil {
		t.Error(err)
	}

	// Create assembler
	asm := &assembler.Assembler{
		Tail:  tail,
		Cache: cache,
	}

	// Assemble package
	// defer os.RemoveAll("/tmp/asm")
	if err := asm.Assemble(key, "/tmp/asm"); err != nil {
		t.Error(err)
	}
}
