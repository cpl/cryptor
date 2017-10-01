package assembler

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

func TestAssembler(t *testing.T) {
	t.Parallel()

	// Create temporary dir for test
	tmpDir, err := ioutil.TempDir("/tmp", "assembler")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create temp cache
	cache, err := cachedb.NewLDBCache(tmpDir, 0, 0)
	if err != nil {
		t.Error(err)
	}

	// Generate random data
	var buffer bytes.Buffer
	buffer.Write(crypt.RandomData(10740))

	// Create chunker
	c := &chunker.Chunker{
		Size:   1024,
		Cache:  cache,
		Reader: &buffer,
	}

	// Chunk data
	tail, err := c.Chunk(crypt.NullKey)
	if err != nil {
		t.Error(err)
	}

	// Create assembler
	asm := &Assembler{
		Tail:  tail,
		Cache: cache,
	}

	// Start assembling package
	err = asm.Assemble(crypt.NullKey)
	if err != nil {
		t.Error(err)
	}
}
