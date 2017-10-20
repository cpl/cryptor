package cachedb

// ManagerConfig ...
type ManagerConfig struct {
	MaxCacheSize  int
	MaxChunkCount int
	MaxChunkSize  int
	MinChunkSize  int

	CurrentCacheSize  int
	CurrentChunkCount int
}
