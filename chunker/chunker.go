package chunker

import (
	"bytes"
	"io"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/cache"
	"github.com/thee-engineer/cryptor/crypt"
)

// Chunker ...
type Chunker struct {
	Size   uint32
	Reader io.Reader
}

// Chunk ...
func (c *Chunker) Chunk(tKey crypt.AESKey) (pHash []byte, err error) {
	// Count chunks
	count := 0

	// chunkHeader stores information about the chunk
	chunkHeader := NewChunkHeader()
	// chunkContent contains chunker.Size bytes of data from the archive
	chunkContent := make([]byte, c.Size)
	// chunkData contains both the header and the content
	chunkData := new(bytes.Buffer)

	pKey := crypt.NullKey
	pHash = make([]byte, 32)

	for {
		// Read archive content into chunks
		read, err := c.Reader.Read(chunkContent)

		// Check for EOF
		if read == 0 || err == io.EOF {
			break
		}

		// Check for errors
		if err != nil {
			return nil, err
		}

		// Add random padding if needed
		if read < int(c.Size) {
			chunkContent = append(
				chunkContent[:read],
				crypt.RandomData(uint(c.Size)-uint(read))...)
			chunkHeader.Padd = c.Size - uint32(read)
		} else {
			chunkHeader.Padd = 0
		}

		// Compute content hash for future checks
		chunkHeader.Hash = crypt.SHA256Data(chunkContent).Sum(nil)

		// Store previous encryption key inside this chunk's header
		chunkHeader.NKey = pKey

		// Store previous encrypted chunk hash inside this chunk's header
		chunkHeader.Next = pHash

		// Create chunk with header and content
		chunkData.Write(chunkHeader.Bytes())
		chunkData.Write(chunkContent)

		// Generatea a new encryption key for each chunk
		if read < int(c.Size) {
			// Use tail key for the last chunk
			pKey = tKey
		} else {
			pKey = crypt.NewKey()
		}

		// Encrypt chunk
		eData, err := crypt.Encrypt(pKey, chunkData.Bytes())
		if err != nil {
			return nil, err
		}

		// Hash encrypted content
		eHash := crypt.SHA256Data(eData).Sum(nil)

		// Create chunk file
		chunkFile, err := os.Create(
			path.Join(cache.GetCachePath(), string(crypt.Encode(eHash))))
		if err != nil {
			return nil, err
		}
		defer chunkFile.Close()

		// Write encrypted data to chunk file
		_, err = chunkFile.Write(eData)
		if err != nil {
			return nil, err
		}

		// Update previous hash
		pHash = eHash

		// Reset buffer
		chunkData.Reset()
		count++
	}

	return pHash, nil
}
