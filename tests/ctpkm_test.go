package tests

import (
	"testing"

	"github.com/thee-engineer/cryptor/ctpkm"
)

func TestPackageCreation(t *testing.T) {
	ctpkm.NewPackageFromDir(".", "testPkg")
}
