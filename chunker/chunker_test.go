package chunker_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

func TestChunker(t *testing.T) {
	t.Parallel()

	// Simulate data
	var buffer bytes.Buffer
	buffer.Write(crypt.RandomData(7747))

	// Create temporary dir for chunks
	tmpDir, err := ioutil.TempDir("/tmp", "cryptor")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create cache for test
	cache, err := ldbcache.NewLDBCache(tmpDir, 0, 0)
	if err != nil {
		t.Error(err)
	}

	// Start chunking data
	if _, err := chunker.ChunkFrom(&buffer, 1024, cache, aes.NullKey); err != nil {
		t.Error(err)
	}
}
