package chunker

import (
	"io"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/cache"
	"github.com/thee-engineer/cryptor/crypt"
)

// Chunker ...
type Chunker struct {
	Size   uint32
	Cache  string
	Reader io.Reader
}

// Chunk ...
func (c Chunker) Chunk(tKey crypt.AESKey) (pHash []byte, err error) {
	// Count chunks
	var count int

	// Check for chunk cache directory
	cache.CheckPath(cache.CacheDir)

	// Make a chunk structure
	chunk := NewChunk(c.Size)

	// Prepare previous hash and key
	pKey := crypt.NullKey
	pHash = make([]byte, 32)

	for {
		// Read archive content into chunks
		read, err := c.Reader.Read(chunk.Content)

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
			chunk.Content = append(
				chunk.Content[:read],
				crypt.RandomData(uint(c.Size)-uint(read))...)
			chunk.Header.Padd = c.Size - uint32(read)
		} else {
			chunk.Header.Padd = 0
		}

		// Compute content hash for future checks
		chunk.Header.Hash = crypt.SHA256Data(chunk.Content[:read]).Sum(nil)

		// Store previous encryption key inside this chunk's header
		chunk.Header.NKey = pKey

		// Store previous encrypted chunk hash inside this chunk's header
		chunk.Header.Next = pHash

		// Generatea a new encryption key for each chunk
		if read < int(c.Size) {
			// Use tail key for the last chunk
			pKey = tKey
		} else {
			pKey = crypt.NewKey()
		}

		// Encrypt chunk data
		eData, err := crypt.Encrypt(pKey, chunk.Bytes())
		if err != nil {
			return nil, err
		}

		// Hash encrypted content
		eHash := crypt.SHA256Data(eData).Sum(nil)

		// Create chunk file
		chunkFile, err := os.Create(
			path.Join(c.Cache, string(crypt.Encode(eHash))))
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

		// Count chunks
		count++
	}

	return pHash, nil
}
