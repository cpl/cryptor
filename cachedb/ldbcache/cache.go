package ldbcache

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/thee-engineer/cryptor/cachedb"
)

// Cache is a LevelDB database used as cache
type Cache struct {
	file string
	db   *leveldb.DB
}

const (
	minCache   = 16
	minHandles = 16
)

// New creates a new LevelDB cache at the given path using a
// number of caches and handles (with a minimum). Returns a Database interface.
func New(file string, cache, handles int) (cachedb.Database, error) {
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

	return &Cache{
		file: file,
		db:   db,
	}, nil
}

// Close read/write access
func (cdb *Cache) Close() error {
	return cdb.db.Close()
}

// Path returns the DB file path
func (cdb *Cache) Path() string {
	return cdb.file
}

// DB returns the LevelDB DB as defined in syndtr/goleveldb/leveldb
// func (cdb *LDBCache) DB() *leveldb.DB {
// 	return cdb.db
// }

// Put stores the key/value pair in the DB
func (cdb *Cache) Put(key, value []byte) error {
	return cdb.db.Put(key, value, nil)
}

// Get returns the value at the given key
func (cdb *Cache) Get(key []byte) ([]byte, error) {
	return cdb.db.Get(key, nil)
}

// Has returns true if the key exists, false otherwise
func (cdb *Cache) Has(key []byte) (bool, error) {
	return cdb.db.Has(key, nil)
}

// Del removes key/value pair from given key
func (cdb *Cache) Del(key []byte) error {
	return cdb.db.Delete(key, nil)
}

// NewBatch returns a new Batch for a LevelDB cache
func (cdb *Cache) NewBatch() cachedb.Batch {
	return &Batch{
		db:    cdb.db,
		batch: new(leveldb.Batch),
	}
}

// NewIterator return a new LevelDB iterator
func (cdb *Cache) NewIterator() cachedb.Iterator {
	return cdb.db.NewIterator(nil, nil)
}
