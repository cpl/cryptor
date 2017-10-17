package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mkideal/cli"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

// Arguments for chunker command
type chunkArg struct {
	cli.Helper
	Source string `cli:"*source" usage:"provide source file/dir for chunker"`
	Size   int    `cli:"size" usage:"chunk size, default 1MB" dft:"1048576"`
	List   bool   `cli:"!list" usage:"lists all chunks in the cache"`
	Pass   string `pw:"p,password" usage:"password for encryption" prompt:"password"`
	Cache  string `cli:"*c,cache" usage:"provide cache source" dft:"default"`
}

// Command for chunker
var chunkCLI = &cli.Command{
	Name: "chunk",
	Desc: "begin chunking the source using the given chunk size and cache",
	Argv: func() interface{} { return new(chunkArg) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*chunkArg)

		// Check which cache to use
		cachePath := ""
		switch argv.Cache {
		case "default":
			cachePath = cachedb.GetCryptorDir()
		case "temp":
			cachePath, _ = ioutil.TempDir("/tmp", "cryptor_cache")
		default:
			cachePath = argv.Cache
		}

		// Open the cache
		cache, err := cachedb.NewLDBCache(cachePath, 16, 16)
		handleErr(err)

		// List chunks and exit
		if argv.List {
			os.Exit(0)
		}

		// Create buffer for creating .tar.gz
		var buffer bytes.Buffer
		err = archive.TarGz(argv.Source, &buffer)
		handleErr(err)

		// Create chunker
		c := &chunker.Chunker{
			Size:   uint32(argv.Size),
			Cache:  cache,
			Reader: &buffer,
		}

		// Derive key from password and start chunking
		hash, err := c.Chunk(aes.NewKeyFromPassword(argv.Pass))
		handleErr(err)

		// Print the tail hash for the final chunk
		fmt.Printf("tail hash: %s\n", crypt.EncodeString(hash))
		return nil
	},
}
