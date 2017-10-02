package crypt

import (
	"crypto/rand"
	"io"
)

// RandomData returns a []byte of given size, containg random data
func RandomData(size uint) []byte {
	data := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err)
	}
	return data
}
