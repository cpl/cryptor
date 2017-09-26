package tests

import (
	"fmt"
	"testing"

	"github.com/thee-engineer/cryptor/assembler"
	"github.com/thee-engineer/cryptor/crypt"
)

func TestAssembler(t *testing.T) {
	eChunk := assembler.GetEChunk("fd8f4599b1c4f725cc68bd40d71f4c40c2716b508674986025c76b94434cd688")
	chunk := eChunk.Decrypt(crypt.NewKeyFromString("ed06859fe6a56fcacca926222b8dd34064494c79ccd71c67f8340f6db0cc8c3a"))

	for !chunk.IsLast() {
		eChunk = assembler.GetEChunk(crypt.EncodeString(chunk.Header.Next))
		chunk = eChunk.Decrypt(chunk.Header.NKey)

		fmt.Println(crypt.EncodeString(chunk.Header.Hash))
	}
}
