package main

import (
	"fmt"

	"github.com/thee-engineer/cryptor/crypto"
)

func main() {
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
