package crypt

import (
	"crypto/rand"
	"io"
)

// RandomData ...
func RandomData(size uint) []byte {
	data := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err)
	}
	return data
}
