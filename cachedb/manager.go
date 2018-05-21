package cachedb

import (
	"log"

	"github.com/thee-engineer/cryptor/common/paths"
)

type manager struct {
	db    Database // Underlaying cache database
	path  string   // Location of cachedb
	size  uint     // Size in bytes
	count uint     // Number of chunks
}

// New ...
func New(path string, db Database) Manager {
	m := &manager{
		db:    db,
		size:  0,
		count: 0,
		path:  path,
	}

	// Update size
	if err := m.update(); err != nil {
		log.Panic(err)
	}

	// Update count
	iter := m.db.NewIterator()
	for iter.Next() {
		m.count++
	}
	iter.Release()

	return m
}

func (m *manager) update() error {
	size, err := paths.DirSize(m.path)
	if err != nil {
		return err
	}
	m.size = size
	return nil
}

// Size ...
func (m *manager) Size() uint {
	size, err := paths.DirSize(m.path)
	if err != nil {
		return m.size
	}
	m.size = size
	return size
}

// Count ...
func (m *manager) Count() uint {
	var count uint
	iter := m.db.NewIterator()
	for iter.Next() {
		count++
	}
	iter.Release()
	m.count = count
	return count
}

// Add ...
func (m *manager) Add(data []byte) error {
	if err := m.db.Put(data); err != nil {
		return err
	}
	// go m.update()
	m.count++
	return nil
}

// Has ...
func (m *manager) Has(key []byte) bool {
	if val, err := m.db.Has(key); !val || err != nil {
		return false
	}
	return true
}

// Get ...
func (m *manager) Get(key []byte) ([]byte, error) {
	return m.db.Get(key)
}

// Del ...
func (m *manager) Del(key []byte) error {
	if err := m.db.Del(key); err != nil {
		return err
	}
	// go m.update()
	m.count--
	return nil
}

// Iterator ...
func (m *manager) Iterator() Iterator {
	return m.db.NewIterator()
}

// Close ...
func (m *manager) Close() error {
	return m.db.Close()
}
