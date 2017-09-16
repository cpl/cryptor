package chunker

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	tmpDirPath    = "/tmp"
	tmpDirPref    = "cryptor"
	chunkFilePath = "%s/%d_chunk"
	chunkSize     = 64
)

// ChunkData ...
func ChunkData(data []byte) (count uint64, tmpDirName string, err error) {

	if len(data) < chunkSize {
		return 0, "", &chunkerError{"Chunker data size too small", 100}
	}

	// Compress entire file
	cDataBuffer, err := compress(data)
	if err != nil {
		return 0, "", err
	}

	// Create tmp directory for chunks
	tmpDirName, err = ioutil.TempDir(tmpDirPath, tmpDirPref)
	if err != nil {
		return 0, "", err
	}

	// Split gzip data into chunks and write to files
	for {
		// Create chunk byte araray of chunkSize
		chunk := make([]byte, chunkSize)

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
		chunkFile := fmt.Sprintf(chunkFilePath, tmpDirName, count)
		err = ioutil.WriteFile(chunkFile, chunk[:read], 0400)
		if err != nil {
			return 0, "", err
		}

		// Count chunks created
		count++
	}

	profile := PackageProfile{"hash", "testpkg", uint64(0), chunkSize, count}
	profile.Generate(tmpDirName)

	return count, tmpDirName, nil
}

// ChunkFile ...
func ChunkFile(filePath string) (count uint64, tmpDirName string, err error) {
	// Read file content
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, "", err
	}

	// Chunk file contents
	return ChunkData(fileContent)
}

// ChunkStdin ...
func ChunkStdin() (count uint64, tmpDirName string, err error) {
	// Read stdin content
	stdinContent, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return 0, "", err
	}

	// Chunk stdin contents
	return ChunkData(stdinContent)
}
