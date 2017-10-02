package cachedb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// LDBCache is a LevelDB database used as cache
type LDBCache struct {
	file string
	db   *leveldb.DB
}

const (
	minCache   = 16
	minHandles = 16
)

// NewLDBCache creates a new LevelDB cache at the given path using a
// number of caches and handles (with a minimum). Returns a Database interface.
func NewLDBCache(file string, cache, handles int) (Database, error) {
	// Check for minimum caches and handles
	if cache < minCache {
		cache = minCache
	}
	if handles < minHandles {
		handles = minHandles
	}

	// Open/Create DB, if file does not exist a new DB will be created
	db, err := leveldb.OpenFile(file, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB,
		Filter:                 filter.NewBloomFilter(10),
	})
	if err != nil {
		return nil, err
	}

	return &LDBCache{
		file: file,
		db:   db,
	}, nil
}

// Close read/write access
func (cdb *LDBCache) Close() error {
	return cdb.db.Close()
}

// Path returns the DB file path
func (cdb *LDBCache) Path() string {
	return cdb.file
}

// DB returns the LevelDB DB as defined in syndtr/goleveldb/leveldb
func (cdb *LDBCache) DB() *leveldb.DB {
	return cdb.db
}

// Put stores the key/value pair in the DB
func (cdb *LDBCache) Put(key, value []byte) error {
	return cdb.db.Put(key, value, nil)
}

// Get returns the value at the given key
func (cdb *LDBCache) Get(key []byte) ([]byte, error) {
	return cdb.db.Get(key, nil)
}

// Has returns true if the key exists, false otherwise
func (cdb *LDBCache) Has(key []byte) (bool, error) {
	return cdb.db.Has(key, nil)
}

// Del removes key/value pair from given key
func (cdb *LDBCache) Del(key []byte) error {
	return cdb.db.Delete(key, nil)
}
