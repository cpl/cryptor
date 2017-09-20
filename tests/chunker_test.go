package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypto"
)

const (
	inputFileName  = "chunker_test.dat"
	outputFileName = "chunker_test.tar.gz"
)

func TestChunking(t *testing.T) {
	// Create input file
	file, err := os.Create(inputFileName)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	// Fill input file with random data, 20KB
	file.Write(crypto.RandomData(20000))

	// Create output file for tar.gz
	outFile, err := os.Create(outputFileName)
	if err != nil {
		t.Error(err)
	}
	defer outFile.Close()

	// Compress file
	if err := archive.TarGz(inputFileName, outFile); err != nil {
		t.Error(err)
	}

	file.Close()
	outFile.Close()

	count, path, err := chunker.ChunkFile(outputFileName, 1024)
	if err != nil {
		fmt.Println(err)
		t.Error(nil)
	}
	fmt.Printf("Chunks: %d @ %s\n", count, path)
}
