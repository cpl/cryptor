package ldbcache_test

import (
	"bytes"
	"os"
	"testing"
)

func TestLDBatch(t *testing.T) {
	t.Parallel()

	// Prepare test environment
	var testData = []string{"", "world", "1409", "\x00cd16\x00", ""}
	dbPath, cache, err := createTestEnv()
	if err != nil {
		t.Error(err)
	}
	defer cache.Close()
	defer os.RemoveAll(dbPath)

	// Create batch and start putting
	batch := cache.NewBatch()
	for _, data := range testData {
		batch.Put([]byte(data), []byte(data))
	}
	batch.Write()

	// Check data
	for _, data := range testData {
		value, err := cache.Get([]byte(data))
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(value, []byte(data)) {
			t.Error("value error: unexpected value")
		}
	}

	// Delete in batch
	for _, data := range testData {
		batch.Del([]byte(data))
	}
	batch.Write()

	// Check data
	for _, data := range testData {
		_, err := cache.Get([]byte(data))
		if err == nil {
			t.Error("value error: got deleted value")
		}
	}
}
