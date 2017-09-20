package tests

import (
	"fmt"
	"testing"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypto"
)

func TestChunkerHead(t *testing.T) {

	header := chunker.NewChunkHeader()
	header.Hash = crypto.SHA256Data(crypto.RandomData(32)).Sum(nil)
	header.Padd = 1000
	header.NKey = crypto.NewKey()[:]

	fmt.Println(len(header.Bytes()))
	fmt.Println(header.Bytes())
	fmt.Println(string(crypto.Encode(header.Bytes())))
}
