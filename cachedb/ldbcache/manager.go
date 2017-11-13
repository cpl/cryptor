package ldbcache

import (
	"errors"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// Manager ...
type Manager struct {
	cachedb.ManagerConfig
	DB cachedb.Database
}

// NewManager ...
func NewManager(config cachedb.ManagerConfig, db cachedb.Database) cachedb.Manager {
	if !cachedb.ValidateConfig(config) {
		panic(cachedb.ErrInvalidConfig)
	}

	ldbMan := &Manager{config, db}
	ldbMan.update()

	return ldbMan
}

func (man *Manager) update() {
	// Open DB iterator
	iter := man.DB.NewIterator()
	defer iter.Release()

	// Count total chunk size and count
	var size int
	var count int

	for iter.Next() {
		count++
		size += len(iter.Value())
	}

	// Check for exceding limits
	if count > man.MaxChunkCount {
		panic(errors.New("ldb man: chunk count excedes max count"))
	}
	if size > man.MaxCacheSize {
		panic(errors.New("ldb man: cache size excedes max size"))
	}

	// Update size
	man.CurrentCacheSize = size
	man.CurrentChunkCount = count
}

func (man *Manager) updateCount() {
	// Open DB iterator
	iter := man.DB.NewIterator()

	// Count entries
	var count int
	for iter.Next() {
		count++
	}

	// Check for exceding limit
	if count > man.MaxChunkCount {
		panic(errors.New("ldb man: chunk count excedes max count"))
	}

	// Update count
	man.CurrentChunkCount = count
}

// Count ...
func (man *Manager) Count() int {
	return man.CurrentChunkCount
}

// Size ...
func (man *Manager) Size() int {
	return man.CurrentCacheSize
}

// Add ...
func (man *Manager) Add(data []byte) error {
	// Ensure chunk count is within limits
	if man.CurrentChunkCount >= man.MaxChunkCount {
		return errors.New("ldb man: chunk limit reached")
	}

	// Ensure chunk is not too large
	chunkSize := len(data)
	if chunkSize < man.MinChunkSize || man.MaxChunkSize < chunkSize {
		return errors.New("ldb man: chunk size outside range limit")
	}

	// Ensure cache has enough space
	if man.CurrentCacheSize+chunkSize > man.MaxCacheSize {
		return errors.New("ldb man: can't add new chunk, will exceed limit")
	}

	// Compute hash as key and add chunk to cache
	if err := man.DB.Put(hashing.SHA256Digest(data), data); err != nil {
		return err
	}

	// Update chunk count and cache size
	man.CurrentChunkCount++
	man.CurrentCacheSize += chunkSize

	return nil
}

// Has ...
func (man *Manager) Has(hex string) bool {
	key, err := crypt.DecodeString(hex)
	if err != nil {
		return false
	}

	has, err := man.DB.Has(key)
	if err != nil {
		return false
	}

	return has
}

// Get ...
func (man *Manager) Get(hex string) ([]byte, error) {
	// Decode key, also validate key
	key, err := crypt.DecodeString(hex)
	if err != nil {
		return nil, err
	}

	// Check that key exists
	if !man.Has(hex) {
		return nil, errors.New("man: no entry with key")
	}

	return man.DB.Get(key)
}

// Del ...
func (man *Manager) Del(hex string) error {
	if man.CurrentChunkCount == 0 {
		man.updateCount()
	}
	if man.CurrentChunkCount == 0 {
		return errors.New("man: cache is empty")
	}

	// Decode key, also validate key
	key, err := crypt.DecodeString(hex)
	if err != nil {
		return err
	}

	data, err := man.Get(hex)
	if err != nil {
		return errors.New("man: failed to find entry with key")
	}

	// Remove data
	if err := man.DB.Del(key); err != nil {
		return err
	}

	// Update count and size
	man.CurrentChunkCount--
	man.CurrentCacheSize -= len(data)

	return nil
}
