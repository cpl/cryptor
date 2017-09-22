package main

import (
	"fmt"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/ctpkg"
	"github.com/thee-engineer/cryptor/utility"
)

// HelpMsg ...
const HelpMsg = `Usage: cryptor <file/dir> <name>`

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	utility.CheckPath(utility.CryptorCachePath)
	utility.CheckPath(utility.CryptorPacksPath)

	if len(os.Args) < 3 {
		fmt.Println(HelpMsg)
		return
	}

	pkg := ctpkg.NewCTPKG(path.Join(cwd, os.Args[1]), os.Args[2], 1024, nil)

	err = pkg.Save()
	if err != nil {
		panic(err)
	}
}
