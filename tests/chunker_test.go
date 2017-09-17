package tests

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/chunker"
)

const (
	testSmallFile = "chunker_test.txt"
	testLargeFile = "random_test_data.txt"
	testDataSize  = 1000000
)

func TestChunkFile(t *testing.T) {
	// Chunk a small file
	_, _, err := chunker.ChunkFile(testSmallFile)
	if err != chunker.ErrorDataSize {
		t.Errorf("Expected ErrorDataSize, got %s", err)
	}

	// Chunk a big random file
	count, path, err := chunker.ChunkFile(testLargeFile)
	if err != nil {
		t.Error(err)
	}

	// Count chunk files
	chunkFiles, err := ioutil.ReadDir(path)
	if err != nil {
		t.Error(err)
	}

	// Check the chunk counts
	if uint(len(chunkFiles)) != count {
		t.Errorf("Mismatch chunkFile count and chunk count")
	}

	// Remove test chunk files
	os.RemoveAll(path)
}

func TestChunkData(t *testing.T) {
	// Generate random test data
	data := make([]byte, testDataSize)
	io.ReadFull(rand.Reader, data)

	// Chunk data
	count, path, err := chunker.ChunkData(data)
	if err != nil {
		t.Error(err)
	}

	// Count chunk files
	chunkFiles, err := ioutil.ReadDir(path)
	if err != nil {
		t.Error(err)
	}

	// Check the chunk counts
	if uint(len(chunkFiles)) != count {
		t.Errorf("Mismatch chunkFile count and chunk count")
	}

	// Assemble data back
	aData, err := chunker.AssembleData(path)
	if err != nil {
		t.Error(err)
	}

	// Check for data mismatch
	if bytes.Compare(data, aData) != 0 {
		t.Error("Mismatch original data and assembled data")
	}

	// Remove test chunks
	os.RemoveAll(path)
}
