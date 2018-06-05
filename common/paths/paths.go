package paths

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
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

// DirSize ...
func DirSize(path string) (uint, error) {
	var size uint
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += uint(info.Size())
		}
		return err
	})
	return size, err
}
