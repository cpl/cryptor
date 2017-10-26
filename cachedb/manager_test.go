package cachedb_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
)

func TestLDBManagerConfig(t *testing.T) {
	t.Parallel()

	validConf := cachedb.ManagerConfig{
		MaxCacheSize:  1000000,
		MaxChunkSize:  10000,
		MinChunkSize:  0,
		MaxChunkCount: 0,
	}

	confs := []cachedb.ManagerConfig{
		// MaxCacheSize < MaxChunkSize
		cachedb.ManagerConfig{
			MaxCacheSize:  100,
			MaxChunkSize:  200,
			MinChunkSize:  10,
			MaxChunkCount: 10,
		},
		// MaxChunkSize < MinChunkSize
		cachedb.ManagerConfig{
			MaxCacheSize:  100,
			MaxChunkSize:  40,
			MinChunkSize:  1000,
			MaxChunkCount: 10,
		},
		// MaxCacheSize < 0
		cachedb.ManagerConfig{
			MaxCacheSize:  -10,
			MaxChunkSize:  40,
			MinChunkSize:  30,
			MaxChunkCount: 10,
		},
		// MaxChunkCount < 0
		cachedb.ManagerConfig{
			MaxCacheSize:  100,
			MaxChunkSize:  40,
			MinChunkSize:  30,
			MaxChunkCount: -10,
		},
		// MaxChunkSize < 0
		cachedb.ManagerConfig{
			MaxCacheSize:  100,
			MaxChunkSize:  -10,
			MinChunkSize:  50,
			MaxChunkCount: 10,
		},
		// MinChunkSize < 0
		cachedb.ManagerConfig{
			MaxCacheSize:  100,
			MaxChunkSize:  40,
			MinChunkSize:  -10,
			MaxChunkCount: 10,
		},
		// MinChunkSize < MaxCacheSize
		cachedb.ManagerConfig{
			MaxCacheSize:  100,
			MaxChunkSize:  40,
			MinChunkSize:  10,
			MaxChunkCount: 100,
		},
		// MaxChunkSize < MinChunkSize
		cachedb.ManagerConfig{
			MaxCacheSize:  1000,
			MaxChunkSize:  10,
			MinChunkSize:  20,
			MaxChunkCount: 10,
		},
	}

	for index, conf := range confs {
		if cachedb.ValidateConfig(conf) {
			t.Errorf("man config: parsed invalid config %d", index)
		}
	}

	if !cachedb.ValidateConfig(validConf) {
		t.Error("man config: ignored valid config")
	}
}
