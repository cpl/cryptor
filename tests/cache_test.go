package tests

import (
	"testing"

	"github.com/thee-engineer/cryptor/cache"
)

func TestCache(t *testing.T) {
	cache.ListChunks()
	cache.ListPacks()
	cache.ClearCache()
	cache.ClearPacks()
}
