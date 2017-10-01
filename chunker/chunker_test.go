package chunker

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
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
	cache, err := cachedb.NewLDBCache(tmpDir, 0, 0)
	if err != nil {
		t.Error(err)
	}

	// Create chunker
	chunker := &Chunker{
		Size:   1024,
		Cache:  cache,
		Reader: &buffer,
	}

	// Start chunking the data
	if _, err := chunker.Chunk(crypt.NullKey); err != nil {
		t.Error(err)
	}

	// // Get chunks file info
	// chunks, err := ioutil.ReadDir(tmpDir)
	// if err != nil {
	// 	panic(err)
	// }

	// // Validate chunk sizes
	// for _, chunk := range chunks {
	// 	if int(chunk.Size()) != 1152 {
	// 		t.Errorf("Invalid chunk size, got %d expected %d",
	// 			chunk.Size(), 1152)
	// 	}
	// }
}
