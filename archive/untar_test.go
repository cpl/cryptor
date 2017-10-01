package archive

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestUnTarEmptyInput(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	// Check for empty input, EOF
	if err := UnTarGz("data/out/out.dat", &buffer); err != nil {
		if err != io.EOF {
			t.Error(err)
		}
	}
}

func TestUnTarWrongInput(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	buffer.Write([]byte{10, 20, 30, 40, 50, 60, 70})

	// Check for unexpected EOF
	if err := UnTarGz("data/out/out.dat", &buffer); err != nil {
		if err.Error() != "unexpected EOF" {
			t.Error(err)
		}
	}

	buffer.Write(crypt.RandomData(200))

	// Check for invalid header
	if err := UnTarGz("data/out/out.dat", &buffer); err != nil {
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

	os.Remove("data/out/noout.dat")

	// Get a valid tar archive
	if err := TarGz("data/tarfile.txt", &buffer); err != nil {
		t.Error(err)
	}

	// Check for empty input, EOF
	if err := UnTarGz("data/out/noout.dat", &buffer); err != nil {
		if err.Error() != "open data/out/noout.dat: is a directory" {
			t.Error(err)
		}
		return
	}

	t.Fail()
}
