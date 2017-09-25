package assembler

import (
	"fmt"

	"github.com/thee-engineer/cryptor/crypt"
)

// Assemble ...
func Assemble(hash string, key *crypt.AESKey) {
	chunk, err := GetChunk(hash, key)
	if err != nil {
		panic(err)
	}

	fmt.Println(chunk.Content)
}
