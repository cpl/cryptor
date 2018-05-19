// Package cachedb offers interfaces for acessing multiple database types. A
// database must have key/pair architecture and must support batch writes.
package cachedb

// Database is an interface that covers read/write interactions with
// underlaying databases.
type Database interface {
	Put(value []byte) error         // Store key/value pair in the DB
	Get(key []byte) ([]byte, error) // Get value with key from the DB
	Has(key []byte) (bool, error)   // Check if key exists in DB
	Del(key []byte) error           // Delete value with key
	Close() error                   // Close I/O with DB
	Path() string                   // Returns the DB file path
	NewBatch() Batch                // Create a Batch for the DB
	NewIterator() Iterator          // Return a new iterator for the DB
}

// Batch contains a set of write instructions that can be executed by a
// Database.
type Batch interface {
	Put(value []byte) error // Store key/value pair instruction
	Del(key []byte) error   // Delete value with key instruction
	Write() error           // Write all instructions in the Batch
}

// Iterator allows iterating over all elements in the DB in a order.
type Iterator interface {
	Next() bool       // Set the iterator to the next item
	Prev() bool       // Set the iterator to the previous item
	Key() []byte      // Returns the key of the current item
	Value() []byte    // Returns the value of the current item
	Seek([]byte) bool // Set the iteratorto the key[]byte param
	Release()         // Release the iterator from use
}

// Manager provides operations on top of the cryptor cache
type Manager interface {
	Size() uint  // Returns the size (in bytes) of the current cache
	Count() uint // Return the total count of chunks

	Add([]byte) error           // Add a new chunk to the cache
	Has([]byte) bool            // Check if cache has chunk
	Get([]byte) ([]byte, error) // Get a chunk from the cache
	Del([]byte) error           // Remove a chunk

	Close() error // Closes the underlaying cache
}
