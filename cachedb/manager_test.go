package cachedb_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/common/con"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

func TestManager(t *testing.T) {
	t.Parallel()

	tmpDir, err := ioutil.TempDir("/tmp", "cachedb_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cache, err := ldbcache.New(tmpDir, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()

	manager := cachedb.NewManager(tmpDir, cache)
	if manager.Count() != 0 {
		t.Errorf("found too many chunks")
	}

	manager.Add(crypt.RandomData(10 * con.MB))
	manager.Add(crypt.RandomData(con.MB))
	manager.Add(crypt.RandomData(2 * con.MB))

	data := crypt.RandomData(con.KB)
	hash := hashing.Hash(data)
	if manager.Has(hash) {
		t.Errorf("found invalid key, Has()")
	}
	if _, err := manager.Get(hash); err == nil {
		t.Errorf("found invalid key, Get()")
	}
	if err := manager.Del(hash); err != nil {
		t.Errorf("found invalid key, Del()")
	}
	manager.Add(data)
	if !manager.Has(hash) {
		t.Errorf("could not find valid key, Has()")
	}
	if _, err := manager.Get(hash); err != nil {
		t.Errorf("could not find valid key, Get()")
		t.Error(err)
	}
	if err := manager.Del(hash); err != nil {
		t.Error(err)
	}
	if err := manager.Close(); err != nil {
		t.Error(err)
	}

	reCache, err := ldbcache.New(tmpDir, 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer reCache.Close()

	newManager := cachedb.NewManager(tmpDir, reCache)
	count := newManager.Count()
	size := newManager.Size()
	if count != 3 {
		t.Errorf("found %d chunks, expected %d", count, 3)
	}
	if size < 13661000 {
		t.Errorf("size is %d, expected more than %d", size, 13661000)
	}
}
