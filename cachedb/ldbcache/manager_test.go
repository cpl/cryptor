package ldbcache_test

import (
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

const testPath = "/tmp/cryptor_man_test_"

func testInt(name string, expected, real int, t *testing.T) {
	if expected != real {
		t.Errorf("ldb man: expecting %s %d, got %d", name, expected, real)
	}
}

var validTestData = []string{
	"hello world",
	"db123",
	"        ",
	"error 482 someone shot the server with a 12 gauge",
}

func TestLDBManager(t *testing.T) {
	t.Parallel()

	db, err := ldbcache.NewLDBCache(testPath+"d", 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testPath + "d")

	conf := cachedb.ManagerConfig{
		MaxCacheSize:  5120,
		MaxChunkSize:  1024,
		MinChunkSize:  5,
		MaxChunkCount: 10,
	}

	man := ldbcache.NewManager(conf, db)
	testInt("count", 0, man.Count(), t)
	testInt("size", 0, man.Size(), t)

	if err := man.Add([]byte("hello world")); err != nil {
		t.Error(err)
	}
	testInt("count", 1, man.Count(), t)
	testInt("size", 11, man.Size(), t)

	if err := man.Del(hashing.SHA256HexDigest([]byte("hello world"))); err != nil {
		t.Error(err)
	}
}

func TestLDBCacheMultiAdd(t *testing.T) {
	t.Parallel()

	db, err := ldbcache.NewLDBCache(testPath+"md", 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testPath + "md")

	conf := cachedb.ManagerConfig{
		MaxCacheSize:  5120,
		MaxChunkSize:  1024,
		MinChunkSize:  5,
		MaxChunkCount: 10,
	}

	man := ldbcache.NewManager(conf, db)

	expectedSize := 0
	expectedCount := len(validTestData)

	for _, testString := range validTestData {
		data := []byte(testString)

		if err := man.Add([]byte(data)); err != nil {
			t.Error(err)
			t.FailNow()
		}

		expectedSize += len(data)
	}

	testInt("count", expectedCount, man.Count(), t)
	testInt("size", expectedSize, man.Size(), t)
}
