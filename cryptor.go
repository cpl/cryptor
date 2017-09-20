package main

import (
	"bytes"
	"io/ioutil"

	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

func main() {
	f, err := ioutil.ReadFile("LICENSE")
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(f)

	chunker := &chunker.Chunker{
		Size:   1024,
		Reader: buf,
		Key:    crypt.NewKey(),
	}

	err = chunker.Start()
	if err != nil {
		panic(err)
	}
}
