package main

import "github.com/thee-engineer/cryptor/ctpkg"

func main() {
	pkg := ctpkg.LoadCTPKG("919a2593aa6c635edc21fdb15c519a0961ac4a66d870e085637c71e4a16844f7")
	pkg.Assemble()
}
