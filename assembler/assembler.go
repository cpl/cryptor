// Package assembler contains data structures and functions that
// read, decrypt and reassembel the original data files from the chunks.
package assembler

import (
	"bytes"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
)

// Assembler contains the tail chunk hash in bytes and the cache containing
// the encrypted chunks.
type Assembler struct {
	Tail  []byte
	Cache cachedb.Database
}

// Assemble starts decrypting the tail chunk with the given AES Key. The
// process extracts the next chunk's data from the current header. If a chunk
// is not found during the assembly process, a network request will be sent
// to the known peers, asking for the missing chunk.
func (a *Assembler) Assemble(key crypt.AESKey) error {
	var cBuffer bytes.Buffer // Chunk buffer, content (no header)
	var aBuffer bytes.Buffer // Assembly buffer, final package

	// Request chunk from cache
	eChunk := GetEChunk(a.Tail, a.Cache)
	// Decrypt given chunk with given key
	chunk, err := eChunk.Decrypt(key)
	if err != nil {
		return err
	}

	// Store decrypted chunk (including header)
	cBuffer.Write(chunk.Content)

	// Process chunks until a final chunk is passed
	for !chunk.IsLast() {
		// Get the next chunk by using the header.Next hash
		eChunk = GetEChunk(chunk.Header.Next, a.Cache)
		// Decrypt the next chunk
		chunk, err = eChunk.Decrypt(chunk.Header.NKey)
		if err != nil {
			return err
		}

		// Store chunk content
		cBuffer.Write(chunk.Content)
	}

	chunkSize := len(chunk.Content)
	bufferLen := cBuffer.Len()
	bufferData := cBuffer.Bytes()

	// Walk trough all processed chunks, place the chunks in the right order
	// inside the assembly buffer
	for index := bufferLen; index > chunkSize; index -= chunkSize {
		aBuffer.Write(bufferData[index-chunkSize : index])
	}

	// Write the final chunk content
	aBuffer.Write(bufferData[0 : bufferLen%chunkSize])

	// Start extracting the .tar.gz archive
	archive.UnTarGz("untar", &aBuffer)

	return nil
}
