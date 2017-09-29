package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/cache"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/ctpkg"
)

func main() {

	sKey := flag.String("key", "", "AES256 Key with hex encoding")
	name := flag.String("name", "", "Optional package name")
	size := flag.Int("size", 1048576, "Chunk size in bytes")
	file := flag.String("file", ".", "Source file/dir to chunk")
	save := flag.Bool("save", false, "If true, a package json file will be created")

	flag.Parse()

	// Check that cryptor packs and chunk cache dir exist
	cache.CheckPath(cache.CryptorCachePath)
	cache.CheckPath(cache.CryptorPacksPath)

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Create Package info
	pkg, err := ctpkg.NewCTPKG(
		path.Join(cwd, *file), // Source
		*name,                         // Name
		uint32(*size),                 // Chunk Size
		crypt.NewKeyFromString(*sKey)) // PKey
	if err != nil {
		panic(err)
	}

	// Store package
	if *save {
		err = pkg.Save()
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("tail: %s\ntkey: %s\n", pkg.Tail, pkg.TKey)
}
