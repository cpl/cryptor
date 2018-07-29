package ldbcache_test

import (
	"math/rand"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

const dataSize = 500000

func BenchmarkLDBWriteMin(b *testing.B) {

	rand.Seed(1)

	cache, err := ldbcache.New("/tmp/cryptor_ldb_bench", 0, 0)
	defer os.RemoveAll("/tmp/cryptor_ldb_bench")
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()

	for count := 0; count < b.N; count++ {
		if err := cache.Put(crypt.RandomData(dataSize)); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkLDBWrite2048(b *testing.B) {

	rand.Seed(1)

	cache, err := ldbcache.New("/tmp/cryptor_ldb_bench", 2048, 0)
	defer os.RemoveAll("/tmp/cryptor_ldb_bench")
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()

	for count := 0; count < b.N; count++ {
		if err := cache.Put(crypt.RandomData(dataSize)); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkLDBWriteSingle(b *testing.B) {

	rand.Seed(1)

	cache, err := ldbcache.New("/tmp/cryptor_ldb_bench", 0, 0)
	defer os.RemoveAll("/tmp/cryptor_ldb_bench")
	if err != nil {
		b.Error(err)
	}

	data := crypt.RandomData(dataSize)

	b.ResetTimer()

	for count := 0; count < b.N; count++ {
		if err := cache.Put(data); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkLDBReadSingle(b *testing.B) {

	rand.Seed(1)

	cache, err := ldbcache.New("/tmp/cryptor_ldb_bench", 0, 0)
	defer os.RemoveAll("/tmp/cryptor_ldb_bench")
	if err != nil {
		b.Error(err)
	}

	data := crypt.RandomData(dataSize)
	key := hashing.Hash(data)
	cache.Put(data)

	b.ResetTimer()

	for count := 0; count < b.N; count++ {
		if _, err := cache.Get(key); err != nil {
			b.Error(err)
		}
	}
}
