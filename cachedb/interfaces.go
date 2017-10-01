package cachedb

// Database ...
type Database interface {
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Del(key []byte) error
	Close() error
	Path() string
	NewBatch() Batch
}

// Batch ...
type Batch interface {
	Put(key, value []byte) error
	Del(key []byte) error
	Write() error
}
