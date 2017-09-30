package cache

import (
	"os"
	"os/user"
	"path"
)

// CryptorPath ...
const CryptorPath = ".cryptor"

// PacksDir ...
const PacksDir = "packs"

// CacheDir ...
const CacheDir = "chunks"

// AssemblyDir ...
const AssemblyDir = "assembly"

// CheckPath ...
func CheckPath(dirPath string) {
	fullPath := path.Join(getUserHomePath(), CryptorPath, dirPath)

	_, err := os.Stat(fullPath)
	if err != nil {
		if err := os.MkdirAll(fullPath, 0700); err != nil {
			panic(err)
		}
	}
}

// GetPacksPath ...
func GetPacksPath() string {
	return path.Join(getUserHomePath(), CryptorPath, PacksDir)
}

// GetCachePath ...
func GetCachePath() string {
	return path.Join(getUserHomePath(), CryptorPath, CacheDir)
}

// GetAssemblyPath ...
func GetAssemblyPath() string {
	return path.Join(getUserHomePath(), CryptorPath, AssemblyDir)
}

func getUserHomePath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}
