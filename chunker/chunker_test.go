package chunker_test

import (
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/chunker"
)

func TestChunker(t *testing.T) {
	t.Parallel()

	db, err := ldbcache.New("/tmp/cryptordb", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("/tmp/cryptordb")
	manager := cachedb.New("/tmp/cryptordb", db)

	c := chunker.New(100, manager)
	archive.TarGz("../test/data/", c)

	_, err = c.Pack("testpassword")
	if err != nil {
		t.Error(err)
	}
}
