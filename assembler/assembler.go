package assembler

import (
	"bytes"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/crypt"
)

// Assembler ...
type Assembler struct {
	Tail  string
	Cache string
}

// Assemble ...
func (a *Assembler) Assemble(key string) error {
	var cBuffer bytes.Buffer
	var aBuffer bytes.Buffer

	eChunk := GetEChunk(a.Tail, a.Cache)
	chunk, err := eChunk.Decrypt(crypt.NewKeyFromString(key))
	if err != nil {
		return err
	}

	cBuffer.Write(chunk.Content)

	for !chunk.IsLast() {
		eChunk = GetEChunk(crypt.EncodeString(chunk.Header.Next), a.Cache)
		chunk, err = eChunk.Decrypt(chunk.Header.NKey)
		if err != nil {
			return err
		}

		cBuffer.Write(chunk.Content)
	}

	chunkSize := len(chunk.Content)
	bufferLen := cBuffer.Len()
	bufferData := cBuffer.Bytes()

	for index := bufferLen; index > chunkSize; index -= chunkSize {
		aBuffer.Write(bufferData[index-chunkSize : index])
	}

	aBuffer.Write(bufferData[0 : bufferLen%chunkSize])
	archive.UnTarGz("untar", &aBuffer)

	return nil
}
