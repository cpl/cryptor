package cache

import (
	"fmt"
	"io/ioutil"
	"os"
)

// ListChunks ...
func ListChunks() {
	CheckPath(CryptorCachePath)
	path := GetCachePath()
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fmt.Println("CHUNKs")

	if len(dirs) == 0 {
		fmt.Println("None")
	}

	for index, dir := range dirs {
		fmt.Printf("%08d | %s\n", index, dir.Name())
	}
	fmt.Println()
}

// ListPacks ...
func ListPacks() {
	CheckPath(CryptorPacksPath)
	path := GetPacksPath()
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fmt.Println("CTPKGs")

	if len(dirs) == 0 {
		fmt.Println("None")
	}

	for index, dir := range dirs {
		fmt.Printf("%04d | %s\n", index, dir.Name())
	}
	fmt.Println()
}

// ClearCache ...
func ClearCache() {
	os.RemoveAll(GetCachePath())
	CheckPath(CryptorCachePath)
}

// ClearPacks ...
func ClearPacks() {
	os.RemoveAll(GetPacksPath())
	CheckPath(CryptorPacksPath)
}
