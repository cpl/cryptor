package ldbcache_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/crypt"
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

	testData := []byte("hello world")
	testInvalidHashes := []string{
		"",
		"aabbcc",
		"INVALID",
		"InVaLiiD",
	}

	// Create new cache
	db, err := ldbcache.NewLDBCache(testPath+"d", 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testPath + "d")

	// Prepare manager config
	conf := cachedb.ManagerConfig{
		MaxCacheSize:  5120,
		MaxChunkSize:  1024,
		MinChunkSize:  5,
		MaxChunkCount: 10,
	}

	// Create manager with config and cache
	man := ldbcache.NewManager(conf, db)
	testInt("count", 0, man.Count(), t)
	testInt("size", 0, man.Size(), t)

	// Add data
	testHash := hashing.SHA256HexDigest(testData)
	if err := man.Add(testData); err != nil {
		t.Error(err)
	}

	// Test Has
	if status := man.Has(testHash); !status {
		t.Error("man: failed to find entry")
	}

	// Get data from cache
	value, err := man.Get(testHash)
	if err != nil {
		t.Error(err)
	}

	// Compare data with data from cache
	if !bytes.Equal(value, testData) {
		t.Error("man: data mismatch")
	}

	// Test size and count
	testInt("count", 1, man.Count(), t)
	testInt("size", 11, man.Size(), t)

	// Delete entry
	if err := man.Del(testHash); err != nil {
		t.Error(err)
	}

	// Test Has (with no data)
	if status := man.Has(testHash); status {
		t.Error("man: found deleted entry")
	}

	// Get data from cache (with no data)
	value, err = man.Get(testHash)
	if err == nil {
		t.Error("man: got deleted entry")
	}
	if value != nil {
		t.Error("man: got non nil data")
	}

	for _, hash := range testInvalidHashes {
		// Invalid key as arg for Has
		if status := man.Has(hash); status {
			t.Error("man: found invalid key")
		}

		// Invalid key as arg for Del
		if err := man.Del(hash); err == nil {
			t.Errorf("man: deleted invalid entry, %s", hash)
		}

		// Invalid key as arg for Get
		value, err = man.Get(hash)
		if err == nil {
			t.Error("man: got invalid entry")
		}
		if value != nil {
			t.Error("man: got invalid data")
		}
	}
}

func TestLDBCacheMultiAdd(t *testing.T) {
	t.Parallel()

	// Create cache
	db, err := ldbcache.NewLDBCache(testPath+"md", 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testPath + "md")

	// Manager config
	conf := cachedb.ManagerConfig{
		MaxCacheSize:  5120,
		MaxChunkSize:  1024,
		MinChunkSize:  5,
		MaxChunkCount: 10,
	}

	// New manager
	man := ldbcache.NewManager(conf, db)

	// Expected test results
	expectedSize := 0
	expectedCount := len(validTestData)

	// Add data to cache using manager
	for _, testString := range validTestData {
		data := []byte(testString)

		if err := man.Add([]byte(data)); err != nil {
			t.Error(err)
			t.FailNow()
		}

		expectedSize += len(data)
	}

	// Compare results
	testInt("count", expectedCount, man.Count(), t)
	testInt("size", expectedSize, man.Size(), t)
}

func TestLDBLimitsSize(t *testing.T) {
	t.Parallel()

	// Create new cache
	db, err := ldbcache.NewLDBCache(testPath+"lims", 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.RemoveAll(testPath + "lims")

	// Prepare manager config
	conf := cachedb.ManagerConfig{
		MaxCacheSize:  100,
		MaxChunkSize:  90,
		MinChunkSize:  5,
		MaxChunkCount: 5,
	}
	if !cachedb.ValidateConfig(conf) {
		t.Error(cachedb.ErrInvalidConfig)
	}

	// Create manager with config and cache
	man := ldbcache.NewManager(conf, db)

	// Test chunk size limits
	if err := man.Add(crypt.RandomData(1)); err == nil {
		t.Error("man: added chunk too small")
	}
	if err := man.Add(crypt.RandomData(91)); err == nil {
		t.Error("man: added chunk too big")
	}

	// Test size and count
	testInt("count", 0, man.Count(), t)
	testInt("size", 0, man.Size(), t)

	// Add max chunk size
	if err := man.Add(crypt.RandomData(90)); err != nil {
		t.Error(err)
	}
	// Add min chunk size
	if err := man.Add(crypt.RandomData(5)); err != nil {
		t.Error(err)
	}

	// Try to exceed max cache size
	if err := man.Add(crypt.RandomData(10)); err == nil {
		t.Error("man: added chunk too big, breaks max cache size")
	}

	// Fill cache size to max
	if err := man.Add(crypt.RandomData(5)); err != nil {
		t.Error(err)
	}

	// Try to exceed max cache size, with min chunk size
	if err := man.Add(crypt.RandomData(5)); err == nil {
		t.Error("man: added to full cache")
	}
}

func TestLDBLimitsCount(t *testing.T) {
	t.Parallel()

	// Create new cache
	db, err := ldbcache.NewLDBCache(testPath+"limc", 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.RemoveAll(testPath + "limc")

	// Prepare manager config
	conf := cachedb.ManagerConfig{
		MaxCacheSize:  100,
		MaxChunkSize:  90,
		MinChunkSize:  5,
		MaxChunkCount: 5,
	}
	if !cachedb.ValidateConfig(conf) {
		t.Error(cachedb.ErrInvalidConfig)
	}

	// Create manager with config and cache
	man := ldbcache.NewManager(conf, db)

	// Add 10 "chunks"
	for count := 0; count < conf.MaxChunkCount; count++ {
		if err := man.Add(crypt.RandomData(10)); err != nil {
			t.Error(err)
		}
	}

	if err := man.Add(crypt.RandomData(10)); err == nil {
		t.Error("man: added to full cache (chunks)")
	}
}
