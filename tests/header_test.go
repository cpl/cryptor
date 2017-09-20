package tests

import (
	"fmt"
	"testing"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

func TestChunkerHead(t *testing.T) {

	header := chunker.NewChunkHeader()
	header.Hash = crypt.SHA256Data(crypt.RandomData(32)).Sum(nil)
	header.Padd = 1000
	header.NKey = crypt.NewKey()[:]

	fmt.Println(len(header.Bytes()))
	fmt.Println(header.Bytes())
	fmt.Println(string(crypt.Encode(header.Bytes())))
}
