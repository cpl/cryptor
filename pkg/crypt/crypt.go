package crypt

import (
	"crypto/rand"
	"encoding/binary"
)

// ZeroBytes takes at least one byte array, iterates the arrays, and sets each byte to zero.
func ZeroBytes(data ...[]byte) {
	for _, set := range data {
		for index := range set {
			set[index] = 0
		}
	}
}

// RandomBytes generates a byte array of given size, containing random byte
// values extracted from "crypto/rand". If you wish to fill an array with random
// data, simply call `rand.Read(arr[:])`.
func RandomBytes(size uint) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}

// RandomUint64 generate and return a random unsigned 64 bit (8 bytes) integer using `crypto/rand`.
func RandomUint64() uint64 {
	return binary.LittleEndian.Uint64(RandomBytes(8))
}
