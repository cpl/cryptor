package ldbcache

import (
	"errors"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// LDBManager ...
type LDBManager struct {
	cachedb.ManagerConfig
	DB cachedb.Database
}

// NewManager ...
func NewManager(config cachedb.ManagerConfig, db cachedb.Database) cachedb.DBManager {
	if !cachedb.ValidateConfig(config) {
		panic(cachedb.ErrInvalidConfig)
	}

	return &LDBManager{config, db}
}

func (man *LDBManager) update() {
	man.updateCount()
	man.updateSize()
}

func (man *LDBManager) updateCount() {
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

func (man *LDBManager) updateSize() {
	// Open DB iterator
	iter := man.DB.NewIterator()

	// Count individual chunk size
	var size int
	for iter.Next() {
		size += len(iter.Value())
	}

	// Check for exceding limit
	if size > man.MaxCacheSize {
		panic(errors.New("ldb man: cache size excedes max size"))
	}

	// Update size
	man.CurrentCacheSize = size
}

// Count ...
func (man *LDBManager) Count() int {
	return man.CurrentChunkCount
}

// Size ...
func (man *LDBManager) Size() int {
	return man.CurrentCacheSize
}

// Add ...
func (man *LDBManager) Add(data []byte) error {
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
func (man *LDBManager) Has(hex string) bool {
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
func (man *LDBManager) Get(hex string) ([]byte, error) {
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
func (man *LDBManager) Del(hex string) error {
	// Decode key, also validate key
	key, err := crypt.DecodeString(hex)
	if err != nil {
		return err
	}

	// Check that key exists
	if !man.Has(hex) {
		return errors.New("man: no entry with key")
	}

	return man.DB.Del(key)
}
