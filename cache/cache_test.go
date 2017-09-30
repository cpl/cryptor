package cache

import "testing"

func TestCache(t *testing.T) {
	getUserHomePath()
	GetCachePath()
	GetAssemblyPath()
	GetPacksPath()
	CheckPath(CacheDir)
	ListChunks()
	ListPacks()
}
