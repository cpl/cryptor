package tests

import (
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

func TestChunker(t *testing.T) {
	f, err := os.Open("header_test.go")
	if err != nil {
		t.Error(err)
	}

	chunker := chunker.NewChunker(f, 32, crypt.NewKey())
	err = chunker.Start()
	if err != nil {
		t.Error(err)
	}
}
