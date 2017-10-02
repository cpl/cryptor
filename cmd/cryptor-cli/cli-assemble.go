package main

import (
	"io/ioutil"

	"github.com/mkideal/cli"
	"github.com/thee-engineer/cryptor/assembler"
	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
)

// Arguments for assembler command
type assembleArg struct {
	cli.Helper
	Hash  string `cli:"*hash" usage:"tail hash, point of start for assembly"`
	Pass  string `pw:"p,password" usage:"password for encryption" prompt:"password"`
	Cache string `cli:"*c,cache" usage:"provide cache source" dft:"default"`
}

// Command for assembler
var assembleCLI = &cli.Command{
	Name: "assemble",
	Desc: "assemble chunks by providing tail hash and the key",
	Argv: func() interface{} { return new(assembleArg) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*assembleArg)

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

		// Open the default cache
		cache, err := cachedb.NewLDBCache(cachePath, 16, 16)
		handleErr(err)

		// Decode tail hash
		tailHash, err := crypt.DecodeString(argv.Hash)
		handleErr(err)

		// Prepare assembler
		a := &assembler.Assembler{
			Tail:  tailHash,
			Cache: cache,
		}

		// Start assembly
		err = a.Assemble(crypt.NewKeyFromPassword(argv.Pass))
		handleErr(err)

		return nil
	},
}
