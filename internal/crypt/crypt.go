package crypt

import (
	"crypto/rand"
	"encoding/binary"
)

// ZeroBytes will iterate each given array and set each byte to 0.
func ZeroBytes(data ...[]byte) {
	for _, set := range data {
		for index := range set {
			set[index] = 0
		}
	}
}

// RandomBytes will return a byte array containing `size` random bytes.
func RandomBytes(size uint) []byte {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		panic(err)
	}
	return data
}

// RandomUint64 will generate 8 random bytes and return them as a uint64.
func RandomUint64() uint64 {
	return binary.LittleEndian.Uint64(RandomBytes(8))
}
