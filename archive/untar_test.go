package archive_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/crypt"
)

func TestUnTarEmptyInput(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	// Check for empty input, EOF
	if err := archive.UnTarGz("../test/data/out/out.dat", &buffer); err != nil {
		if err != io.EOF {
			t.Error(err)
		}
	}
}

func TestUnTarWrongInput(t *testing.T) {
	// t.Parallel()

	var buffer bytes.Buffer

	buffer.Write([]byte{10, 20, 30, 40, 50, 60, 70})

	// Check for unexpected EOF
	if err := archive.UnTarGz("../test/data/out/out.dat", &buffer); err != nil {
		if err.Error() != "unexpected EOF" {
			t.Error(err)
		}
	}

	buffer.Write(crypt.RandomData(200))

	// Check for invalid header
	if err := archive.UnTarGz("../test/data/out/out.dat", &buffer); err != nil {
		if err.Error() != "gzip: invalid header" {
			t.Error(err)
		}
		return
	}

	t.Fail()
}

func TestUnTarNoOutputFile(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	// Get a valid tar archive
	if err := archive.TarGz("../test/data/tarfile.txt", &buffer); err != nil {
		t.Error(err)
	}

	// Check for empty output, EOF
	defer os.Remove("../test/data/out/noout.dat")
	if err := archive.UnTarGz("../test/data/out/noout.dat", &buffer); err != nil {
		if err != nil {
			t.Error(err)
		}
	}

}

func TestUnTarNoOutputDir(t *testing.T) {
	// t.Parallel()

	var buffer bytes.Buffer

	// Get a valid tar archive
	if err := archive.TarGz("../test/data/", &buffer); err != nil {
		t.Error(err)
	}

	// Check for empty output, EOF
	defer os.RemoveAll("tmpdat")
	if err := archive.UnTarGz("tmpdat/down/the/test", &buffer); err != nil {
		if err != nil {
			t.Error(err)
		}
	}

}
