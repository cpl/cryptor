package crypt

import (
	"crypto/rand"
	"io"
)

// RandomData returns a []byte of given size, containg random data
func RandomData(size uint) []byte {
	data := make([]byte, size)
	io.ReadFull(rand.Reader, data)
	return data
}
