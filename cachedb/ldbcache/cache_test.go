package ldbcache_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/common/paths"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/encode/b16"
	"github.com/thee-engineer/cryptor/crypt/hashing"
	"github.com/thee-engineer/cryptor/utils"
)

// Test data for all key/value pair tests
var testData = []string{"", "world", "1409", "\x00cd16\x00", ""}

func createTestEnv() (string, cachedb.Database, error) {
	// Create tmp dir for test
	tmpDir, err := ioutil.TempDir("/tmp", "cachedb_test")
	if err != nil {
		return "", nil, err
	}

	// Create test db
	// Use 0 cache and 0 handlers, should default to min
	cdb, err := ldbcache.New(tmpDir, 0, 0)
	if err != nil {
		return "", nil, err
	}

	return tmpDir, cdb, nil
}

func TestCDBBasic(t *testing.T) {
	t.Parallel()

	// Create test env
	dbPath, cdb, err := createTestEnv()
	defer os.RemoveAll(dbPath)
	utils.CheckErrTest(err, t)

	defer cdb.Close()

	// Check db path
	if cdb.Path() != dbPath {
		t.Error("wrong path")
	}

	// Test Put
	tdata := []byte("hello")
	if err := cdb.Put(tdata); err != nil {
		t.Error(err)
	}

	// Test Has (true)
	status, err := cdb.Has(hashing.Hash(tdata))
	utils.CheckErrTest(err, t)

	if !status {
		t.Error("key error: expected key not found")
	}

	// Test Has (false)
	status, err = cdb.Has([]byte("item1"))
	utils.CheckErrTest(err, t)

	if status {
		t.Error("key error: found unexpected key")
	}

	// Test Del
	if err := cdb.Del([]byte("item0")); err != nil {
		t.Error(err)
	}

	// Check if deleted key exists
	status, err = cdb.Has([]byte("item0"))
	utils.CheckErrTest(err, t)

	if status {
		t.Error("key error: found deleted key")
	}
}

func TestCDBSameKeyPut(t *testing.T) {
	t.Parallel()

	// Create test env
	dbPath, cdb, err := createTestEnv()
	defer os.RemoveAll(dbPath)
	utils.CheckErrTest(err, t)

	defer cdb.Close()

	// Put count entries with the same keys and random data
	for count := 0; count < 5; count++ {
		data := b16.Encode(crypt.RandomData(10))
		cdb.Put(data)
	}

	// TODO: Add iterator to check keys
}

func TestCDBAdvanced(t *testing.T) {
	t.Parallel()

	// Create test env
	dbPath, cdb, err := createTestEnv()
	defer os.RemoveAll(dbPath)
	utils.CheckErrTest(err, t)

	defer cdb.Close()

	// Put
	for _, data := range testData {
		err := cdb.Put([]byte(data))
		utils.CheckErrTest(err, t)
	}

	// Get
	for _, data := range testData {
		value, err := cdb.Get(hashing.Hash([]byte(data)))
		utils.CheckErrTest(err, t)

		if !bytes.Equal(value, []byte(data)) {
			t.Error("value error: got unexpected value")
		}
	}

	// Put overwrite
	// for _, data := range testData {
	// 	err := cdb.Put([]byte(data), []byte("OVERWRITE"))
	// 	utils.CheckErrTest(err, t)
	// }

	// Get overwrite
	// for _, data := range testData {
	// 	value, err := cdb.Get(hashing.Hash([]byte(data)))
	// 	utils.CheckErrTest(err, t)

	// 	if !bytes.Equal(value, []byte("OVERWRITE")) {
	// 		t.Error("value error: got unexpected value")
	// 	}
	// }

	// Del
	for _, data := range testData {
		err := cdb.Del([]byte(data))
		utils.CheckErrTest(err, t)

	}

	// Get del
	for _, data := range testData {
		_, err := cdb.Get([]byte(data))
		if err == nil {
			t.Error("value error: got deleted value")
		}
	}
}

func TestCDBIterator(t *testing.T) {
	t.Parallel()

	// Create test env
	dbPath, cdb, err := createTestEnv()
	defer os.RemoveAll(dbPath)
	utils.CheckErrTest(err, t)

	defer cdb.Close()

	// Put test data in cache
	for _, data := range testData {
		err := cdb.Put([]byte(data))
		utils.CheckErrTest(err, t)

	}

	// Open iterator over cache
	iter := cdb.NewIterator()
	count := 0
	for iter.Next() {
		// Get current key/value pair
		iter.Key()
		iter.Value()
		count++
	}
	if count != len(testData)-1 {
		t.Errorf("iter error: invalid iterator length; got %d; expected %d;",
			count, len(testData)-1)
	}

	iter.Release()
}

func TestCDBErrors(t *testing.T) {
	t.Parallel()

	// Create test env
	dbPath, cdb, err := createTestEnv()
	defer os.RemoveAll(dbPath)
	if err != nil {
		t.Error(err)
	}
	defer cdb.Close()

	// Get (non existing)
	data, err := cdb.Get([]byte("hello"))
	if err == nil {
		t.Error(err)
	}
	if data != nil {
		t.Error("ldb: got invalid data")
	}

	// Has (non existing)
	status, err := cdb.Has([]byte("hello"))
	if status {
		t.Error("ldb: found invalid key/value")
	}

	// Del (non existing)
	if err := cdb.Del([]byte("hello")); err != nil {
		t.Error(err)
	}

	// Put value
	tdata := []byte("world")
	if err := cdb.Put(tdata); err != nil {
		t.Error(err)
	}

	// Del (existing)
	if err := cdb.Del(hashing.Hash(tdata)); err != nil {
		t.Error(err)
	}

	// Get (non existing after delete)
	data, err = cdb.Get([]byte("hello"))
	if err != nil {
		log.Println("err")
	}

	if data != nil {
		if len(data) != 0 {
			t.Error("ldb: got invalid data")
		}
	}

	// Has (non existing after delete)
	status, err = cdb.Has([]byte("hello"))
	utils.CheckErrTest(err, t)

	if status {
		t.Error("ldb: found invalid key/value")
	}

	cdb.Close()

	// Put value
	if err := cdb.Put([]byte("world")); err == nil {
		t.Error("ldb: wrote to closed database")
	}
}

func TestCDBSize(t *testing.T) {
	t.Parallel()

	// Initiate db
	path, _, err := createTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	// Check size of empty ldb
	size, err := paths.DirSize(path)
	if err != nil {
		t.Fatal(err)
	}
	if size == 0 {
		t.Errorf("size is 0")
	}
}
