package main

import (
	"fmt"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"

	"github.com/mkideal/cli"
)

// Arguments for cache command
type cacheArg struct {
	cli.Helper
	Cache string `cli:"*c,cache" usage:"specify cache location" dft:"default"`
}

// Command for cache management
var cacheCLI = &cli.Command{
	Name: "cache",
	Desc: "cache management",
	Argv: func() interface{} { return new(cli.Helper) },
	Fn: func(ctx *cli.Context) error {
		return nil
	},
}

// Command for creating new cache
var newCacheCLI = &cli.Command{
	Name: "add",
	Desc: "cache creation",
	Argv: func() interface{} { return new(cacheArg) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*cacheArg)

		// Create cache
		cache, err := cachedb.NewLDBCache(handleCache(argv.Cache), 0, 0)
		handleErr(err)
		defer cache.Close()

		return nil
	},
}

// Command for listing items in a cache
var listCacheCLI = &cli.Command{
	Name: "list",
	Desc: "cache removal",
	Argv: func() interface{} { return new(cacheArg) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*cacheArg)

		// Open (or create)
		cache, err := cachedb.NewLDBCache(handleCache(argv.Cache), 0, 0)
		handleErr(err)
		defer cache.Close()

		// List items in cache (hash and chunk len + total len)
		iter := cache.NewIterator()
		fmt.Printf("%66s   %16s\n", "chunk hash", "chunk length")
		total := 0
		for iter.Next() {
			key, value := iter.Key(), iter.Value()
			total += len(value)
			fmt.Printf("%66s : %16d\n", crypt.EncodeString(key), len(value))
		}
		fmt.Printf("\n%s : %d B\n", "total", total)
		iter.Release()

		return nil
	},
}
