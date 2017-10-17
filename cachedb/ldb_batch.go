package cachedb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// LDBBatch is a LevelDB batch interface
type LDBBatch struct {
	db    *leveldb.DB
	batch *leveldb.Batch
	size  int
}

// NewBatch returns a new Batch for a LevelDB cache
func (db *LDBCache) NewBatch() Batch {
	return &LDBBatch{
		db:    db.db,
		batch: new(leveldb.Batch),
	}
}

// Put appends a Put key/value pair instruction in the Batch
func (b *LDBBatch) Put(key, value []byte) error {
	b.batch.Put(key, value)
	b.size += len(value)
	return nil
}

// Del appends a Del key instruction in the Batch
func (b *LDBBatch) Del(key []byte) error {
	b.batch.Delete(key)
	return nil
}

// Write applies all instructions inside the Batch on the LevelDB
func (b *LDBBatch) Write() error {
	return b.db.Write(b.batch, nil)
}
