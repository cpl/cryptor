package ctpkg

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestCTPKG(t *testing.T) {
	// Create temporary dir for test
	tmpDir, err := ioutil.TempDir("/tmp", "ctpkg")
	if err != nil {
		t.Error(err)
	}

	// Create a package out of ctpkg.go
	ctpkg, err := NewCTPKG("ctpkg.go", "ctpkg", tmpDir, 1024, crypt.NullKey)
	if err != nil {
		t.Error(err)
	}

	// Intentional error
	if err := ctpkg.Save("fakepath"); err == nil {
		t.Fail()
	}

	// Store CTPKG file
	if err := ctpkg.Save(tmpDir); err != nil {
		t.Error(err)
	}

	// Load CTPKG
	_, err = LoadCTPKG(ctpkg.Tail, tmpDir)
	if err != nil {
		t.Error(err)
	}

	// Remove test files
	os.RemoveAll(tmpDir)
}
