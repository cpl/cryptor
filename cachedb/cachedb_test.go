package cachedb

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func createTestEnv() (string, Database, error) {
	// Create tmp dir for test
	tmpDir, err := ioutil.TempDir("/tmp", "cachedb_test")
	if err != nil {
		return "", nil, err
	}

	// Create test db
	// Use 0 cache and 0 handlers, should default to min
	cdb, err := NewCacheDB(tmpDir, 0, 0)
	if err != nil {
		return "", nil, err
	}

	return tmpDir, cdb, nil
}

func TestCDBBasic(t *testing.T) {
	t.Parallel()

	dbPath, cdb, err := createTestEnv()
	if err != nil {
		t.Error(err)
	}
	defer cdb.Close()

	// Check db path
	if cdb.Path() != dbPath {
		t.Error("Path mismatch")
	}

	// Test Put
	if err := cdb.Put([]byte("item0"), []byte("hello")); err != nil {
		t.Error(err)
	}

	// Test Has (true)
	status, err := cdb.Has([]byte("item0"))
	if err != nil {
		t.Error(err)
	}
	if status != true {
		t.Error("Expected key not found")
	}

	// Test Has (false)
	status, err = cdb.Has([]byte("item1"))
	if err != nil {
		t.Error(err)
	}
	if status != false {
		t.Error("Found unexpected key")
	}

	// Test Del
	if err := cdb.Del([]byte("item0")); err != nil {
		t.Error(err)
	}

	// Check if deleted key exists
	status, err = cdb.Has([]byte("item0"))
	if err != nil {
		t.Error(err)
	}
	if status != false {
		t.Error("Found deleted key")
	}

	// Remove test files
	os.RemoveAll(dbPath)
}

func TestCDBSameKeyPut(t *testing.T) {
	t.Parallel()

	dbPath, cdb, err := createTestEnv()
	if err != nil {
		t.Error(err)
	}
	defer cdb.Close()

	// Put count entries with the same keys and random data
	for count := 0; count < 5; count++ {
		data := crypt.Encode(crypt.RandomData(10))
		cdb.Put([]byte("key"), data)
	}

	os.RemoveAll(dbPath)
}

func TestCDBAdvanced(t *testing.T) {
	t.Parallel()

	dbPath, cdb, err := createTestEnv()
	if err != nil {
		t.Error(err)
	}
	defer cdb.Close()

	var testData = []string{"", "world", "1409", "\x00cd16\x00", ""}

	// Put
	for _, data := range testData {
		err := cdb.Put([]byte(data), []byte(data))
		if err != nil {
			t.Error(err)
		}
	}

	// Get
	for _, data := range testData {
		value, err := cdb.Get([]byte(data))
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(value, []byte(data)) {
			t.Error("Got unexpected value")
		}
	}

	// Put overwrite
	for _, data := range testData {
		err := cdb.Put([]byte(data), []byte("OVERWRITE"))
		if err != nil {
			t.Error(err)
		}
	}

	// Get overwrite
	for _, data := range testData {
		value, err := cdb.Get([]byte(data))
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(value, []byte("OVERWRITE")) {
			t.Error("Got unexpected value")
		}
	}

	// Del
	for _, data := range testData {
		err := cdb.Del([]byte(data))
		if err != nil {
			t.Error(err)
		}
	}

	// Get del
	for _, data := range testData {
		_, err := cdb.Get([]byte(data))
		if err == nil {
			t.Error("Got deleted data")
		}
	}

	os.RemoveAll(dbPath)
}
