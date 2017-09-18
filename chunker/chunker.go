package chunker

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/thee-engineer/cryptor/utility"
)

const (
	tmpDir        = "/tmp"
	tmpDirPrefix  = "cryptor"
	chunkFilePath = "%s/%032d_chunk"

	// ChunkSize ...
	ChunkSize = 1024
)

// ChunkData ...
func ChunkData(data []byte) (count uint, tmpDirPath string, err error) {

	// Check if data fits in  at least two packets
	if len(data) < ChunkSize {
		return count, tmpDirPath, ErrorDataSize
	}

	// Compress entire file
	cDataBuffer, err := utility.Compress(data)
	if err != nil {
		return count, tmpDirPath, err
	}

	// Check if cData fits in  at least two packets
	if cDataBuffer.Len() < ChunkSize {
		return count, tmpDirPath, ErrorDataSizeCompressoin
	}

	// Compute expected chunk count
	expectedChunkCount := uint(cDataBuffer.Len()) / ChunkSize
	if cDataBuffer.Len()%ChunkSize != 0 {
		expectedChunkCount++
	}

	// Create tmp directory for chunks
	tmpDirPath, err = ioutil.TempDir(tmpDir, tmpDirPrefix)
	if err != nil {
		return count, tmpDirPath, err
	}

	// Split gzip data into chunks and write to files
	for {
		// Create chunk byte araray of chunkSize
		chunk := make([]byte, ChunkSize)

		// Read into chunk from gzip data
		read, err := cDataBuffer.Read(chunk)

		// Check EOF
		if read == 0 {
			break
		}

		// Check error
		if err != nil {
			panic(err)
		}

		// Write chunk data to chunk file
		chunkFile := fmt.Sprintf(chunkFilePath, tmpDirPath, count)
		err = ioutil.WriteFile(chunkFile, chunk[:read], 0400)
		if err != nil {
			return 0, "", err
		}

		// Count chunks created
		count++
	}

	if expectedChunkCount != count {
		panic(ErrorChunkCount)
	}

	return count, tmpDirPath, nil
}

// ChunkFile ...
func ChunkFile(filePath string) (count uint, tmpDirPath string, err error) {
	// Read file content
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, "", err
	}

	// Chunk file contents
	return ChunkData(fileContent)
}

// ChunkStdin ...
func ChunkStdin() (count uint, tmpDirPath string, err error) {
	// Read stdin content
	stdinContent, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return 0, "", err
	}

	// Chunk stdin contents
	return ChunkData(stdinContent)
}
