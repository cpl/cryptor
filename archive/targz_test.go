package archive_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/crypt"
)

func TestTarUnTar(t *testing.T) {
	t.Parallel()

	// Hash initial file
	hash, err := crypt.SHA256File("data/tarfile.txt")
	if err != nil {
		t.Error(err)
	}

	// Median buffer
	var buffer bytes.Buffer

	// Tar the file
	if err := archive.TarGz("data/tarfile.txt", &buffer); err != nil {
		t.Error(err)
	}

	// Create output file
	_, err = os.Create("data/out/tarfile.txt")
	if err != nil {
		t.Error(err)
	}

	// Untar the file
	if err := archive.UnTarGz("data/out/tarfile.txt", &buffer); err != nil {
		t.Error(err)
	}

	// Hash file after taring and untaring
	fHash, err := crypt.SHA256File("data/out/tarfile.txt")
	if err != nil {
		t.Error(err)
	}

	// Compare the two hashes
	if !bytes.Equal(fHash.Sum(nil), hash.Sum(nil)) {
		t.Error("Output hash does not match initial hash!")
	}

	// Remove output file
	err = os.Remove("data/out/tarfile.txt")
	if err != nil {
		t.Error(err)
	}
}
