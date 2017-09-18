package tests

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/thee-engineer/cryptor/utility"
)

func TestGzip(t *testing.T) {
	originalData := []byte("Hello, World")
	compData, err := utility.Compress(originalData)
	if err != nil {
		t.Error(err)
	}
	uncmData, err := utility.Decompress(compData)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(originalData, uncmData) != 0 {
		t.Errorf("Data mistmatch")
	}
}

func TestGzipBig(t *testing.T) {
	originalData := make([]byte, 1000000)
	io.ReadFull(rand.Reader, originalData)
	compData, err := utility.Compress(originalData)
	if err != nil {
		t.Error(err)
	}
	uncmData, err := utility.Decompress(compData)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(originalData, uncmData) != 0 {
		t.Errorf("Data mistmatch")
	}
}
