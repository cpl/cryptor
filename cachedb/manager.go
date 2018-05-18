package cachedb

import "errors"

// ManagerConfig ...
type ManagerConfig struct {
	MaxCacheSize  int `json:"maxCacheSize"`
	MaxChunkCount int `json:"maxChunkCount"`
	MaxChunkSize  int `json:"maxChunkSize"`
	MinChunkSize  int `json:"minChunkSize"`

	CurrentCacheSize  int `json:"-"`
	CurrentChunkCount int `json:"-"`
}

// ErrInvalidConfig is used in panic and error returns.
var ErrInvalidConfig = errors.New("man: invalid configuration")

// ValidateConfig takes a Manager configuration object and validates it.
func ValidateConfig(config ManagerConfig) bool {
	// Check for negatives
	if config.MaxChunkCount < 0 || config.MaxCacheSize < 0 {
		return false
	}
	if config.MaxChunkSize < 0 || config.MinChunkSize < 0 {
		return false
	}

	// Make sure 1 chunk fits at least
	if config.MaxCacheSize < config.MinChunkSize {
		return false
	}
	if config.MaxCacheSize < config.MaxChunkSize {
		return false
	}

	// Make sure n small chunks fit
	if config.MaxCacheSize < config.MinChunkSize*config.MaxChunkCount {
		return false
	}

	// Make sure max chunk size >= than min size
	if config.MinChunkSize > config.MaxChunkSize {
		return false
	}

	return true
}
