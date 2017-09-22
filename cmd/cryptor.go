package main

import (
	"fmt"

	"github.com/thee-engineer/cryptor/ctpkg"
	"github.com/thee-engineer/cryptor/utility"
)

func main() {
	utility.CheckPath(utility.CryptorCachePath)
	utility.CheckPath(utility.CryptorPacksPath)
	pkg := ctpkg.NewCTPKG("cryptor.go", "cgo", 1024, nil)
	fmt.Println(pkg.PKey)
}
