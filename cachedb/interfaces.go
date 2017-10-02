// Package cachedb offers interfaces for acessing multiple database types. A
// database must have key/pair architecture and must support batch writes.
package cachedb

// Database is an interface that covers read/write interactions with
// underlaying databases.
type Database interface {
	Put(key, value []byte) error    // Store key/value pair in the DB
	Get(key []byte) ([]byte, error) // Get value with key from the DB
	Has(key []byte) (bool, error)   // Check if key exists in DB
	Del(key []byte) error           // Delete value with key
	Close() error                   // Close I/O with DB
	Path() string                   // Returns the DB file path
	NewBatch() Batch                // Create a Batch for the DB
}

// Batch contains a set of write instructions that can be executed by a
// Database.
type Batch interface {
	Put(key, value []byte) error // Store key/value pair instruction
	Del(key []byte) error        // Delete value with key insturction
	Write() error                // Write all instructions in the Batch
}
