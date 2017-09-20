package chunker

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/crypto"
)

const (
	tmpDir    = "/tmp"
	tmpPrefix = "cryptor"
	chunkName = "chunk%08d"
)

// ChunkData ...
func ChunkData(data []byte, size int, pkey [32]byte) (count int, tmpDirPath string, err error) {
	buffer := bytes.NewBuffer(data)

	// Check if data fits in  at least two chunks
	if len(data) < size {
		return count, "", errorDataSize
	}

	// Create tmp directory for chunks
	tmpDirPath, err = ioutil.TempDir(tmpDir, tmpPrefix)
	if err != nil {
		return count, tmpDirPath, err
	}

	// Create []byte of chunk size
	chunkData := make([]byte, size)

	// Process the HEAD chunk
	read, err := buffer.Read(chunkData)
	if err != nil {
		panic(err)
	}
	nKey := crypto.NewKey()

	chunkHead := NewChunkHeader()
	chunkHead.Size = size
	chunkHead.PadSize = 0
	chunkHead.NKey = crypto.Encode(nKey[:])
	chunkHead.Hash = crypto.Encode(crypto.SHA256Data(chunkData).Sum(nil))

	for {
		// Read data into the chunk
		read, err := buffer.Read(chunkData)

		// EOF
		if read == 0 || err == io.EOF {
			break
		}

		// Check error
		if err != nil {
			return count, tmpDirPath, err
		}

		// Chunk needs padding
		if read != size {
			for index := read; index < size; index++ {
				chunkData[index] = 0
			}
		}

		// Encrypt chunk data
		// eChunkData, err := crypto.Encrypt(cryptoKey, chunkData)
		// if err != nil {
		// 	panic(err)
		// }

		// Create chunk file
		chunkPath := path.Join(tmpDirPath, fmt.Sprintf(chunkName, count))
		chunkFile, err := os.Create(chunkPath)
		if err != nil {
			return count, tmpDirPath, err
		}
		// Write chunk header to file
		_, err = chunkFile.Write(chunkHead.Bytes())
		if err != nil {
			return count, tmpDirPath, err
		}
		// Write chunk content to file
		_, err = chunkFile.Write(chunkData)
		if err != nil {
			return count, tmpDirPath, err
		}

		// Count chunks
		count++
	}

	return count, tmpDirPath, nil
}

// ChunkFile ...
func ChunkFile(filePath string, chunkSize int) (count int, tmpDirPath string, err error) {
	// Read file content
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, "", err
	}

	// Chunk file contents
	return ChunkData(fileContent, chunkSize)
}
