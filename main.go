package main

import (
	"fmt"

	"github.com/thee-engineer/cryptor/chunker"
)

func main() {
	count, path, err := chunker.ChunkFile("testfile.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%d %s\n", count, path)

	data, err := chunker.AssembleData(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}
