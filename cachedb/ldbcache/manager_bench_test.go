package ldbcache_test

import (
	"math/rand"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/crypt"
)

func BenchmarkLDBManagerWrite(b *testing.B) {

	rand.Seed(1)

	// Create new cache
	db, err := ldbcache.NewLDBCache("/tmp/cryptor_ldbman_bench", 0, 0)
	if err != nil {
		b.Error(err)
	}
	defer os.RemoveAll("/tmp/cryptor_ldbman_bench")

	// Create manager config
	conf := cachedb.ManagerConfig{
		MaxCacheSize:  12000 * b.N,
		MaxChunkSize:  12000,
		MinChunkSize:  1,
		MaxChunkCount: b.N + 2,
	}

	// New manager
	man := ldbcache.NewManager(conf, db)

	b.ResetTimer()

	// Benchmark
	for index := 0; index < b.N; index++ {
		if err := man.Add(crypt.RandomData(10000)); err != nil {
			b.Error(err)
		}
	}
}
