package ldbcache_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/crypt/hashing"
	"github.com/thee-engineer/cryptor/utils"
)

func TestLDBatch(t *testing.T) {
	t.Parallel()

	// Prepare test environment
	var testData = []string{"", "world", "1409", "\x00cd16\x00", ""}
	dbPath, cache, err := createTestEnv()
	defer os.RemoveAll(dbPath)
	utils.CheckErrTest(err, t)
	defer cache.Close()

	// Create batch and start putting
	batch := cache.NewBatch()
	for _, data := range testData {
		batch.Put([]byte(data))
	}
	batch.Write()

	// Check data
	for _, data := range testData {
		value, err := cache.Get(hashing.Hash([]byte(data)))
		utils.CheckErrTest(err, t)

		if !bytes.Equal(value, []byte(data)) {
			t.Error("value error: unexpected value")
		}
	}

	// Delete in batch
	for _, data := range testData {
		batch.Del(hashing.Hash([]byte(data)))
	}
	batch.Write()

	// Check data
	for _, data := range testData {
		_, err := cache.Get(hashing.Hash([]byte(data)))
		if err == nil {
			t.Error("value error: got deleted value")
		}
	}
}
