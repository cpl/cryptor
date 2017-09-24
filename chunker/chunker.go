package chunker

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/utility"
)

// Chunker ...
type Chunker struct {
	Size   uint32
	Reader io.Reader
	Key    crypt.AESKey
}

// Chunk ...
func (chunker *Chunker) Chunk() (err error) {
	// Keep count of chunks
	count := 0

	// Prepare data
	chunkHeader := NewChunkHeader()
	chunkData := new(bytes.Buffer)
	chunkCont := make([]byte, chunker.Size)

	key := crypt.NullKey
	var pHash []byte

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

		// Add random padding if needed
		if read < int(chunker.Size) {
			chunkCont = append(
				chunkCont[read:],
				crypt.RandomData(uint(chunker.Size)-uint(read))...)
			chunkHeader.Padd = chunker.Size - uint32(read)
		} else {
			chunkHeader.Padd = 0
		}

		// Store used key for next key
		chunkHeader.NKey = key

		// Store previous hash as next chunk hash
		chunkHeader.Next = pHash

		// Create chunk with header and content
		chunkData.Write(chunkHeader.Bytes())
		chunkData.Write(chunkCont)

		// Prepare keys for encryption
		key = crypt.NewKey()

		// Encrypt chunk
		eData, err := crypt.Encrypt(key, chunkData.Bytes())
		if err != nil {
			return err
		}

		// Hash encrypted content
		eHash := crypt.SHA256Data(eData).Sum(nil)

		// Create chunk file
		chunkFile, err := os.Create(
			path.Join(utility.GetCachePath(), string(crypt.Encode(eHash))))
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		// Write encrypted data to chunk file
		_, err = chunkFile.Write(eData)
		if err != nil {
			return err
		}

		fmt.Println(string(crypt.Encode(pHash)))

		// Update previous hash
		pHash = eHash

		fmt.Println(string(crypt.Encode(pHash)))

		// Reset buffer
		chunkData.Reset()
		count++
	}

	return nil
}
