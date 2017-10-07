package cachedb_test

import (
	"path"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
)

func TestPaths(t *testing.T) {
	t.Parallel()

	// Test that paths works
	absRoot := cachedb.GetCryptorDir()
	if absRoot != path.Join(cachedb.GetUserHomePath(), cachedb.CryptorDir) {
		t.Error("path error: cryptor paths mismatch")
	}
}
