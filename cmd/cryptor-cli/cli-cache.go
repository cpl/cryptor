package main

import "github.com/mkideal/cli"

// Arguments for cache command
type cacheArg struct {
	cli.Helper
}

// Command for cache management
var cacheCLI = &cli.Command{
	Name: "cache",
	Desc: "cache management",
	Argv: func() interface{} { return new(cacheArg) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*cacheArg)
		ctx.String("%s\n", argv)
		return nil
	},
}
