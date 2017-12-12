package cachedb_test

import (
	"os"
	"path"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
)

func TestPaths(t *testing.T) {
	t.Parallel()

	// Test home path
	envHome := os.Getenv("HOME")
	os.Setenv("HOME", "")
	if os.Getenv("HOME") != "" {
		t.Fail()
	}
	usrHome := cachedb.GetUserHomePath()
	if usrHome != envHome {
		t.Fail()
	}
	os.Setenv("HOME", envHome)

	// Test that paths works
	absRoot := cachedb.GetCryptorDir()
	if absRoot != path.Join(cachedb.GetUserHomePath(), cachedb.CryptorDir) {
		t.Error("path error: cryptor paths mismatch")
	}
}
