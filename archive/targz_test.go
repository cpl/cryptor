package archive_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
)

func TestTarUnTar(t *testing.T) {
	t.Parallel()

	// Median buffer
	var buffer bytes.Buffer

	// Tar the file
	if err := archive.TarGz("data/tarfile.txt", &buffer); err != nil {
		t.Error(err)
	}

	// Create output file
	_, err := os.Create("data/out/tarfile.txt")
	if err != nil {
		t.Error(err)
	}

	// Untar the file
	if err := archive.UnTarGz("data/out/tarfile.txt", &buffer); err != nil {
		t.Error(err)
	}

	// Remove output file
	if err := os.Remove("data/out/tarfile.txt"); err != nil {
		t.Error(err)
	}
}

func TestFullTar(t *testing.T) {
	t.Parallel()
	var buffer bytes.Buffer

	// Archive cryptor package
	if err := archive.TarGz("..", &buffer); err != nil {
		t.Error(err)
	}

	// Get a temporary path
	tmpPath, err := ioutil.TempDir("/tmp", "targztest")
	if err != nil {
		t.Error(err)
	}
	// Remove directory, let UnTar create it
	os.RemoveAll(tmpPath)

	// Untar cryptor package
	defer os.RemoveAll(tmpPath)
	if err := archive.UnTarGz(tmpPath, &buffer); err != nil {
		t.Error(err)
	}
}
