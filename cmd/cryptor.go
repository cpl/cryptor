package main

import (
	"fmt"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/crypt"

	"github.com/thee-engineer/cryptor/ctpkg"
	"github.com/thee-engineer/cryptor/utility"
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
	utility.CheckPath(utility.CryptorCachePath)
	utility.CheckPath(utility.CryptorPacksPath)

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
