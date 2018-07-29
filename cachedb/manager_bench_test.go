package cachedb_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/common/con"
	"github.com/thee-engineer/cryptor/crypt"
)

func createTestEnv() (string, cachedb.Database, cachedb.Manager) {
	tmpDir, err := ioutil.TempDir("/tmp", "cachedb_test")
	if err != nil {
		log.Fatal(err)
	}

	cache, err := ldbcache.New(tmpDir, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	manager := cachedb.New(tmpDir, cache)
	if manager.Count() != 0 {
		log.Fatal("found too many chunks")
	}

	return tmpDir, cache, manager
}

// BenchmarkManagerAdd-4 20 83204106 ns/op 4608998 B/op 102 allocs/op
// BenchmarkManagerAdd-4 20 78841192 ns/op 4618482 B/op 150 allocs/op
// BenchmarkManagerAdd-4 20 84433261 ns/op 4618412 B/op 151 allocs/op
func BenchmarkManagerAdd(b *testing.B) {
	tmpDir, cache, manager := createTestEnv()
	defer os.RemoveAll(tmpDir)
	defer cache.Close()

	b.ResetTimer()

	for count := 0; count < b.N; count++ {
		manager.Add(crypt.RandomData(con.MB))
	}
}
