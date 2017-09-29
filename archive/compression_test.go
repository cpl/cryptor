package archive

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestCompression(t *testing.T) {
	_, err := Compress(crypt.RandomData(100))
	if err != nil {
		t.Error(err)
	}

	_, err = Compress([]byte{})
	if err != nil {
		t.Error(err)
	}
}

func TestDecompression(t *testing.T) {
	initialData := crypt.RandomData(100)
	buffer, err := Compress(initialData)
	if err != nil {
		t.Error(err)
	}

	data, err := Decompress(buffer)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(data, initialData) != 0 {
		t.Error("Initial data and uncompressed data do not match!")
	}
}
