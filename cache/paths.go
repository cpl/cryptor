package cache

import (
	"os"
	"os/user"
	"path"
)

// CryptorPacksPath ...
const CryptorPacksPath = ".cryptor/packs"

// CryptorCachePath ...
const CryptorCachePath = ".cryptor/chunks"

// CryptorAssemblyPath ...
const CryptorAssemblyPath = ".cryptor/assembly"

// CheckPath ...
func CheckPath(dirPath string) {
	fullPath := path.Join(getUserHomePath(), dirPath)

	_, err := os.Stat(fullPath)
	if err != nil {
		if err := os.MkdirAll(fullPath, 0700); err != nil {
			panic(err)
		}
	}
}

// GetPacksPath ...
func GetPacksPath() string {
	return path.Join(getUserHomePath(), CryptorPacksPath)
}

// GetCachePath ...
func GetCachePath() string {
	return path.Join(getUserHomePath(), CryptorCachePath)
}

// GetAssemblyPath ...
func GetAssemblyPath() string {
	return path.Join(getUserHomePath(), CryptorAssemblyPath)
}

func getUserHomePath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}
