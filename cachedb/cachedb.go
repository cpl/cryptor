package cachedb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// CacheDB ...
type CacheDB struct {
	file string
	db   *leveldb.DB
}

const (
	minCache   = 16
	minHandles = 16
)

// NewCacheDB ...
func NewCacheDB(file string, cache, handles int) (Database, error) {
	// Check for minimum caches and handles
	if cache < minCache {
		cache = minCache
	}
	if handles < minHandles {
		handles = minHandles
	}

	// Open/Create DB
	db, err := leveldb.OpenFile(file, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB,
		Filter:                 filter.NewBloomFilter(10),
	})
	if err != nil {
		return nil, err
	}

	return &CacheDB{
		file: file,
		db:   db,
	}, nil
}

// Close ...
func (cdb *CacheDB) Close() error {
	return cdb.db.Close()
}

// Path ...
func (cdb *CacheDB) Path() string {
	return cdb.file
}

// DB ...
func (cdb *CacheDB) DB() *leveldb.DB {
	return cdb.db
}

// Put ...
func (cdb *CacheDB) Put(key, value []byte) error {
	return cdb.db.Put(key, value, nil)
}

// Get ...
func (cdb *CacheDB) Get(key []byte) ([]byte, error) {
	return cdb.db.Get(key, nil)
}

// Has ...
func (cdb *CacheDB) Has(key []byte) (bool, error) {
	return cdb.db.Has(key, nil)
}

// Del ...
func (cdb *CacheDB) Del(key []byte) error {
	return cdb.db.Delete(key, nil)
}
