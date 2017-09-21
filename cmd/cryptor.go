package main

import (
	"fmt"

	"github.com/thee-engineer/cryptor/ctpkg"
	"github.com/thee-engineer/cryptor/utility"
)

func main() {
	utility.CheckCache()
	pkg := ctpkg.NewCTPKG("cryptor.go", "cgo", 1024)
	fmt.Println(pkg.PKey)
}
