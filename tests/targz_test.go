package tests

import (
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
)

const (
	testDirectory = "targz_test"
	testOutTar    = "out.tar.gz"
)

func TestTarGz(t *testing.T) {
	archive.TarGz(testDirectory, os.Stdout)
}
