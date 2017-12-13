package assembler_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/assembler"
	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

func createTestCache() ([]byte, error) {
	// Create cache
	db, err := ldbcache.NewLDBCache("data", 0, 0)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	cache := ldbcache.NewManager(cachedb.DefaultManagerConfig, db)

	// Create archive and data buffers
	var buffer bytes.Buffer
	archive.TarGz("chunk.go", &buffer)

	// Chunk files and get tail hash
	tail, err := chunker.ChunkFrom(&buffer, 16, cache, aes.NullKey)
	if err != nil {
		return nil, err
	}

	return tail, nil
}

func TestAssembler(t *testing.T) {
	// Create test cache
	tail, err := createTestCache()
	defer os.RemoveAll("data")
	if err != nil {
		t.Error(err)
	}

	// Open cache
	cache, err := ldbcache.NewLDBCache("data", 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer cache.Close()

	// Create assembler
	asm := assembler.NewDefaultAssembler(tail, cache)

	// Start assembling package
	defer os.RemoveAll("untar")
	err = asm.Assemble(aes.NullKey, "untar")
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
	db, err := ldbcache.NewLDBCache("/tmp/asmcnktest", 16, 16)
	defer os.RemoveAll("/tmp/asmcnktest")
	if err != nil {
		t.Error(err)
	}
	cache := ldbcache.NewManager(cachedb.DefaultManagerConfig, db)

	// Chunk with derived key
	key := aes.NewKeyFromPassword("testing")
	tail, err := chunker.ChunkFrom(&buffer, 1000000, cache, key)
	if err != nil {
		t.Log("cache count:", cache.Count())
		t.Log("cache size:", cache.Size())
		t.Fatal(err)
	}

	// Create assembler
	asm := assembler.NewDefaultAssembler(tail, db)

	// Assemble package
	defer os.RemoveAll("/tmp/asm")
	if err := asm.Assemble(key, "/tmp/asm"); err != nil {
		t.Error(err)
	}
}
