// cryptor
//
//    Usage: cryptor chunk
//
// -h, --help             display help information
// --source               *provide source file/dir for chunker
// --size[=1048576]       chunk size, default 1MB
// --list                 lists all chunks in the cache
// -p, --password         password for encryption
// -c, --cache[=default]  *provide cache source
//
//    Usage: cryptor assemble
//
// -h, --help             display help information
// --hash                 *tail hash, point of start for assembly
// -p, --password         password for encryption
// -c, --cache[=default]  *provide cache source
package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

var help = cli.HelpCommand("display this help message")

func main() {
	// Tree with all possible subcommands
	if err := cli.Root(root,
		cli.Tree(help),
		cli.Tree(chunkCLI),
		cli.Tree(assembleCLI),
		cli.Tree(cacheCLI,
			cli.Tree(newCacheCLI),
			cli.Tree(listCacheCLI),
		),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Print error and exit program
func handleErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Arguments that cen be used in root
type rootArg struct {
	cli.Helper
}

// The way root treats arguments
var root = &cli.Command{
	Desc: "cryptor, file sharing redefined",
	Argv: func() interface{} { return new(rootArg) },
	Fn: func(ctx *cli.Context) error {
		// argv := ctx.Argv().(*rootArg)
		// ctx.String("%s\n", argv)
		return nil
	},
}
