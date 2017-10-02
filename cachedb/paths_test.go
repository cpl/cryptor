package cachedb

import (
	"path"
	"testing"
)

func TestPaths(t *testing.T) {
	t.Parallel()

	// Test that paths works
	absRoot := GetCryptorDir()
	if absRoot != path.Join(getUserHomePath(), cryptorDir) {
		t.Error("path error: cryptor paths mismatch")
	}
}
