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

// ZeroBytes takes a byte array and replaces each byte with 0
func ZeroBytes(data []byte) {
	for index := range data {
		data[index] = 0
	}
}
