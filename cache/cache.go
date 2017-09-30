package cache

import (
	"fmt"
	"io/ioutil"
	"os"
)

// ListChunks ...
func ListChunks() {
	CheckPath(CacheDir)
	path := GetCachePath()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%64s | %s \n", "HASH", "SIZE")
	for _, file := range files {
		fmt.Printf("%64s | %d \n", file.Name(), file.Size())
	}
	fmt.Println()
}

// ListPacks ...
func ListPacks() {
	CheckPath(PacksDir)
	path := GetPacksPath()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%64s\n", "HASH")
	for _, file := range files {
		fmt.Printf("%64s\n", file.Name())
	}
	fmt.Println()
}

// ClearCache ...
func ClearCache() {
	os.RemoveAll(GetCachePath())
	CheckPath(CacheDir)
}

// ClearPacks ...
func ClearPacks() {
	os.RemoveAll(GetPacksPath())
	CheckPath(PacksDir)
}
