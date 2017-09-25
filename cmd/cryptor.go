package main

import (
	"fmt"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/cache"
	"github.com/thee-engineer/cryptor/crypt"

	"github.com/thee-engineer/cryptor/ctpkg"
)

// HelpMsg ...
const HelpMsg = `Usage: cryptor <file/dir> <name>`

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Check that cryptor packs and chunk cache dir exist
	cache.CheckPath(cache.CryptorCachePath)
	cache.CheckPath(cache.CryptorPacksPath)

	// Check that enough args were given
	if len(os.Args) < 3 {
		fmt.Println(HelpMsg)
		return
	}

	// Create Package info
	pkg := ctpkg.NewCTPKG(
		path.Join(cwd, os.Args[1]), // Source
		os.Args[2],                 // Name
		1024,                       // Chunk Size
		crypt.NullKey)              // PKey

	// Store package
	err = pkg.Save()
	if err != nil {
		panic(err)
	}
}
