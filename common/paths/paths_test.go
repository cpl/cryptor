package paths_test

import (
	"os"
	"path"
	"testing"

	"github.com/thee-engineer/cryptor/common/paths"
)

func TestPaths(t *testing.T) {
	t.Parallel()

	// Test home path
	envHome := os.Getenv("HOME")
	os.Setenv("HOME", "")
	if os.Getenv("HOME") != "" {
		t.Fail()
	}
	usrHome := paths.GetUserHomePath()
	if usrHome != envHome {
		t.Fail()
	}
	os.Setenv("HOME", envHome)

	// Test that paths works
	absRoot := paths.GetCryptorDir()
	if absRoot != path.Join(paths.GetUserHomePath(), paths.CryptorDir) {
		t.Error("path error: cryptor paths mismatch")
	}
}
