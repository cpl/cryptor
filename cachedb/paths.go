package cachedb

import (
	"os"
	"os/user"
	"path"
)

const (
	// CryptorDir is the default name of the Cryptor directory
	CryptorDir = ".cryptor"
)

// GetCryptorDir returns the default cryptor cache dir for the user running
// cryptor.
func GetCryptorDir() string {
	return path.Join(GetUserHomePath(), CryptorDir)
}

// GetUserHomePath returns the current user's home dir abs path
func GetUserHomePath() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
