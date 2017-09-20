package chunker

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/crypt"
)

const (
	tmpPath       = "/tmp"
	tmpPref       = "cryptor"
	chunkFileName = "chunk%08d"
)

// Chunker ...
type Chunker struct {
	Size   uint32
	Key    *[32]byte
	Reader io.Reader
}

// Start ...
func (chunker *Chunker) Start() error {
	// Keep count of chunks
	count := 0

	// Prepare data
	chunkHeader := NewChunkHeader()
	chunkData := new(bytes.Buffer)
	chunkCont := make([]byte, chunker.Size)

	// Create temporary directory
	tmpDir, err := ioutil.TempDir(tmpPath, tmpPref)
	if err != nil {
		return err
	}

	// Prepare keys
	var keyNext *[32]byte
	var keyThis *[32]byte

	for {
		// Read archive content into chunks
		read, err := chunker.Reader.Read(chunkCont)

		// Check for EOF
		if read == 0 || err == io.EOF {
			break
		}

		// Check for errors
		if err != nil {
			return err
		}

		// Switch keys
		if count == 0 {
			keyThis = chunker.Key
		} else {
			keyThis = keyNext
		}
		keyNext = crypt.NewKey()

		// Update chunk header
		chunkHeader.Hash = crypt.SHA256Data(chunkCont).Sum(nil)
		chunkHeader.NKey = keyNext[:]
		chunkHeader.Padd = chunker.Size - uint32(read)

		// Add padding if needed
		if read < int(chunker.Size) {
			for index := read; index < int(chunker.Size); index++ {
				chunkCont[index] = 0
			}
		}

		// Create chunk with header and content
		chunkData.Write(chunkHeader.Bytes())
		chunkData.Write(chunkCont)

		// Encrypt chunk
		eData, err := crypt.Encrypt(keyThis, chunkData.Bytes())
		if err != nil {
			return err
		}

		// Create chunk file
		chunkFile, err := os.Create(
			path.Join(tmpDir, fmt.Sprintf(chunkFileName, count)))
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		// Write encrypted data to chunk file
		_, err = chunkFile.Write(eData)
		if err != nil {
			return err
		}

		// Reset buffer
		chunkData.Reset()
		count++
	}

	return nil
}
