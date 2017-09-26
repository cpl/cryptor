package assembler

import (
	"bytes"
	"fmt"

	"github.com/thee-engineer/cryptor/archive"

	"github.com/thee-engineer/cryptor/crypt"
)

// Assemble ...
func Assemble(tail, key string) {
	var cBuffer bytes.Buffer
	var aBuffer bytes.Buffer

	eChunk := GetEChunk(tail)
	chunk := eChunk.Decrypt(crypt.NewKeyFromString(key))

	if !chunk.IsValid() {
		fmt.Println("Invalid chunk!")
	}

	cBuffer.Write(chunk.Content)

	for !chunk.IsLast() {
		eChunk = GetEChunk(crypt.EncodeString(chunk.Header.Next))
		chunk = eChunk.Decrypt(chunk.Header.NKey)

		if !chunk.IsValid() {
			fmt.Println("Invalid chunk!")
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
	fmt.Println(crypt.EncodeString(crypt.SHA256Data(aBuffer.Bytes()).Sum(nil)))

	archive.UnTarGz("untar", &aBuffer)

}
