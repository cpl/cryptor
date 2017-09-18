package chunker

const (
	tmpDir = "/tmp"
)

/*
// ChunkData ...
func ChunkData(data []byte, chunkSize) (count int, tmpDirPath string, err error) {

	// Check if data fits in  at least two packets
	if len(data) < ChunkSize {
		return count, tmpDirPath, errorDataSize
	}

	// Compute expected chunk count
	expectedChunkCount := len(data) / ChunkSize
	if len(data)%ChunkSize != 0 {
		expectedChunkCount++
	}

	// Create tmp directory for chunks
	tmpDirPath, err = ioutil.TempDir(tmpDir, tmpDirPrefix)
	if err != nil {
		return count, tmpDirPath, err
	}

	// Something must have gone wrong
	if expectedChunkCount != count {
		panic(errorChunkCount)
	}

	return count, tmpDirPath, nil
}

// ChunkFile ...
func ChunkFile(filePath string) (count int, tmpDirPath string, err error) {
	// Read file content
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, "", err
	}

	// Chunk file contents
	return ChunkData(fileContent)
}

// ChunkStdin ...
func ChunkStdin() (count int, tmpDirPath string, err error) {
	// Read stdin content
	stdinContent, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return 0, "", err
	}

	// Chunk stdin contents
	return ChunkData(stdinContent)
}
*/
