package utility

import (
	"os"
	"os/user"
	"path"
)

const (
	cryptorCache = ".cryptor/chunks"
)

// CheckCache ...
func CheckCache() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	fullPath := path.Join(usr.HomeDir, cryptorCache)

	_, err = os.Stat(fullPath)
	if err != nil {
		if err := os.MkdirAll(fullPath, 0700); err != nil {
			panic(err)
		}
	}
}

// GetCachePath ...
func GetCachePath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(usr.HomeDir, cryptorCache)
}
