package tests

import (
	"bytes"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
)

func TestTarGz(t *testing.T) {
	buffer := new(bytes.Buffer)

	// Archive test data as .tar.gz
	if err := archive.TarGz("targz_test", buffer); err != nil {
		t.Error(err)
	}

	// Extract test data from .tar.gz
	if err := archive.UnTarGz("targz_test_out", buffer); err != nil {
		t.Error(err)
	}

	os.RemoveAll("targz_test_out")
}
