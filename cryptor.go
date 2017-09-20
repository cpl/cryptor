package main

import (
	"fmt"

	"github.com/thee-engineer/cryptor/ctpkgm"
)

func main() {
	pkg := ctpkgm.NewCTPKG("LICENSE", "license", 1024)
	fmt.Println(pkg.PKey)
}
