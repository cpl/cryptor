package main

import (
	"os"

	"github.com/thee-engineer/cryptor/archive"
)

func main() {
	err := archive.TarGz("LICENSE", os.Stdout)
	if err != nil {
		panic(err)
	}
}
