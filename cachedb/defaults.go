package cachedb

// DefaultManagerConfig is the default recommended configuration for
// caches in GENERAL!
var DefaultManagerConfig = ManagerConfig{
	MaxCacheSize:      268435456,
	MaxChunkCount:     2048,
	MaxChunkSize:      1048576,
	MinChunkSize:      0,
	CurrentCacheSize:  0,
	CurrentChunkCount: 0,
}
