package cachedb

import (
	"os/user"
	"path"
)

const (
	cryptorDir = ".cryptor"
)

// GetCryptorDir ...
func GetCryptorDir() string {
	return path.Join(getUserHomePath(), cryptorDir)
}

func getUserHomePath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}
