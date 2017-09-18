package main

import (
	"fmt"

	"github.com/thee-engineer/cryptor/ctpkgm"
)

func main() {
	pkg := ctpkgm.NewCTPKG("crypto", "cryptor-crypto")
	fmt.Println(pkg.String())
}
