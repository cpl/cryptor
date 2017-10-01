package cachedb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// CDBBatch ...
type CDBBatch struct {
	db    *leveldb.DB
	batch *leveldb.Batch
	size  int
}

// NewBatch ...
func (cdb *CacheDB) NewBatch() Batch {
	return &CDBBatch{
		db:    cdb.db,
		batch: new(leveldb.Batch),
	}
}

// Put ...
func (b *CDBBatch) Put(key, value []byte) error {
	b.batch.Put(key, value)
	b.size += len(value)
	return nil
}

// Del ...
func (b *CDBBatch) Del(key []byte) error {
	b.batch.Delete(key)
	return nil
}

// Write ...
func (b *CDBBatch) Write() error {
	return b.db.Write(b.batch, nil)
}
