package cachedb

import "github.com/boltdb/bolt"

// BoltCache ...
type BoltCache struct {
	file string
	db   *bolt.DB
}

// NewBoltCache ...
func NewBoltCache(file string) (Database, error) {
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &BoltCache{
		file: file,
		db:   db,
	}, nil
}

// Path ...
func (cache *BoltCache) Path() string {
	return cache.file
}

// DB ...
func (cache *BoltCache) DB() *bolt.DB {
	return cache.db
}

// Close ...
func (cache *BoltCache) Close() error {
	return cache.Close()
}

// Put ...
func (cache *BoltCache) Put(key, value []byte) error {
	return nil
}

// Has ...
func (cache *BoltCache) Has(key []byte) (bool, error) {
	return false, nil
}

// Get ...
func (cache *BoltCache) Get(key []byte) ([]byte, error) {
	return nil, nil
}

// Del ...
func (cache *BoltCache) Del(key []byte) error {
	return nil
}

// NewBatch ...
func (cache *BoltCache) NewBatch() Batch {
	return nil
}
