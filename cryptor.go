package main

import (
	"fmt"

	"github.com/thee-engineer/cryptor/crypto"
)

func main() {
	// count, path, err := chunker.ChunkFile("testfile.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%d %s\n", count, path)

	// data, err := chunker.AssembleData(path)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(data))

	msg := []byte("Hello, World!")
	key := crypto.NewAESKey()

	e, _ := crypto.AESEncrypt(key, msg)
	fmt.Println(string(crypto.Encode(e)))

	d, err := crypto.AESDecrypt(key, e)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(d))
}
