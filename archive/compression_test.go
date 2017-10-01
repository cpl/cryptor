package archive

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestCompression(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	initialData := crypt.RandomData(100)
	buffer, err := Compress(initialData)
	if err != nil {
		t.Error(err)
	}

	data, err := Decompress(buffer)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, initialData) {
		t.Error("Initial data and uncompressed data do not match!")
	}
}
