package cachedb

import "github.com/thee-engineer/cryptor/common/con"

// DefaultManagerConfig is the default recommended configuration for
// caches in GENERAL!
var DefaultManagerConfig = ManagerConfig{
	MaxCacheSize:      1 * con.GB,
	MaxChunkCount:     1000,
	MaxChunkSize:      10 * con.MB,
	MinChunkSize:      0,
	CurrentCacheSize:  0,
	CurrentChunkCount: 0,
}
