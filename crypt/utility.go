package crypt // import "cpl.li/go/cryptor/crypt"

import (
	"crypto/rand"
)

// ZeroBytes takes at least one array of bytes and sets each byte individually
// to zero.
func ZeroBytes(data ...[]byte) {
	// iterate each array
	for _, set := range data {
		// iterate array bytes
		for index := range set {
			set[index] = 0
		}
	}
}

// RandomBytes generates a byte array of given size, containing random byte
// values extracted from "crypto/rand". If you wish to fill an array with random
// data, simply call `rand.Read(arr[:])`.
func RandomBytes(size uint) []byte {
	// allocate byte array
	data := make([]byte, size)

	// read random bytes
	rand.Read(data)

	return data
}
