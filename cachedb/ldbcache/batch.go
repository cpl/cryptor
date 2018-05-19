package ldbcache

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// Batch is a LevelDB batch interface
type Batch struct {
	db    *leveldb.DB
	batch *leveldb.Batch
}

// Put appends a Put key/value pair instruction in the Batch
func (b *Batch) Put(value []byte) error {
	b.batch.Put(hashing.Hash(value), value)
	return nil
}

// Del appends a Del key instruction in the Batch
func (b *Batch) Del(key []byte) error {
	b.batch.Delete(key)
	return nil
}

// Write applies all instructions inside the Batch on the LevelDB
func (b *Batch) Write() error {
	return b.db.Write(b.batch, nil)
}
