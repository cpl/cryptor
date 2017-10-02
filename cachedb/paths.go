package cachedb

import (
	"os/user"
	"path"
)

const (
	cryptorDir = ".cryptor"
)

// GetCryptorDir returns the default cryptor cache dir for the user running
// cryptor.
func GetCryptorDir() string {
	return path.Join(getUserHomePath(), cryptorDir)
}

// Get the current user and return the home dir abs path
func getUserHomePath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}
