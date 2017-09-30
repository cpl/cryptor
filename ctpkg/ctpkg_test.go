package ctpkg

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestCTPKG(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "ctpkg")
	if err != nil {
		t.Error(err)
	}

	ctpkg, err := NewCTPKG("ctpkg.go", "ctpkg", tmpDir, 1024, crypt.NullKey)
	if err != nil {
		t.Error(err)
	}

	ctpkg.Save()

	os.RemoveAll(tmpDir)
}
