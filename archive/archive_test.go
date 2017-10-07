package archive_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
)

func TestTarFile(t *testing.T) {
	t.Parallel()

	// Create output file
	file, err := os.Create("data/out/tarfile.tar.gz")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	// Archive input to output
	if err := archive.TarGz("data/tarfile.txt", file); err != nil {
		t.Error(err)
	}

	// Remove file
	if err := os.Remove("data/out/tarfile.tar.gz"); err != nil {
		t.Error(err)
	}
}

func TestTarDir(t *testing.T) {
	t.Parallel()

	// Create output file
	file, err := os.Create("data/out/tardir.tar.gz")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	// Archive input to output
	if err := archive.TarGz("data/tardir", file); err != nil {
		t.Error(err)
	}

	// Remove file
	if err := os.Remove("data/out/tardir.tar.gz"); err != nil {
		t.Error(err)
	}
}

func TestTarErrors(t *testing.T) {
	t.Parallel()

	// Create output buffer
	var buffer bytes.Buffer

	// Try to archive non existent file
	if err := archive.TarGz("data/nosuchfile.txt", &buffer); err != nil {
		if err.Error() != "lstat data/nosuchfile.txt: no such file or directory" {
			t.Error(err)
		}
	}

	// Try to archive file without permissions
	if err := archive.TarGz("data/tar000.txt", &buffer); err != nil {
		if err.Error() != "open data/tar000.txt: permission denied" {
			t.Error(err)
		}
	}

	// Try to archive empty dir
	if err := archive.TarGz("data/emptydir", &buffer); err != nil {
		t.Error(err)
	}
}
