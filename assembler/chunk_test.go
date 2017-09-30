package assembler

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

func TestEChunk(t *testing.T) {
	// Create temporary dir for test
	tmpDir, err := ioutil.TempDir("/tmp", "assembler")
	if err != nil {
		t.Error(err)
	}

	// Test data
	var buffer bytes.Buffer
	data := crypt.RandomData(520)
	buffer.Write(data)

	// Create chunker
	chunker := &chunker.Chunker{
		Size:   1024,
		Cache:  tmpDir,
		Reader: &buffer,
	}
	chunker.Chunk(crypt.NullKey)

	// Get chunk hash
	chunks, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		t.Error(err)
	}

	// Get the single chunk hash
	chunkHash := chunks[0].Name()

	// Read encrypted chunk
	eChunk := GetEChunk(chunkHash, tmpDir)
	dChunk, err := eChunk.Decrypt(crypt.NullKey)
	if err != nil {
		t.Error(err)
	}

	// Invalid hash
	if !dChunk.IsValid() {
		t.Error("Chunk is not valid!")
	}

	// Chunk should be the tail (as it is the only chunk)
	if !dChunk.IsLast() {
		t.Error("Chunk is not tail!")
	}

	// Compare initial data with data after encryption, storage and decryption
	if bytes.Compare(dChunk.Content, data) != 0 {
		t.Error("Data mismatch!")
	}

	// Remove test files
	os.RemoveAll(tmpDir)
}
