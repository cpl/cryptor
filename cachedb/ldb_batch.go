package cachedb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// LDBBatch ...
type LDBBatch struct {
	db    *leveldb.DB
	batch *leveldb.Batch
	size  int
}

// NewBatch ...
func (db *LDBCache) NewBatch() Batch {
	return &LDBBatch{
		db:    db.db,
		batch: new(leveldb.Batch),
	}
}

// Put ...
func (b *LDBBatch) Put(key, value []byte) error {
	b.batch.Put(key, value)
	b.size += len(value)
	return nil
}

// Del ...
func (b *LDBBatch) Del(key []byte) error {
	b.batch.Delete(key)
	return nil
}

// Write ...
func (b *LDBBatch) Write() error {
	return b.db.Write(b.batch, nil)
}
