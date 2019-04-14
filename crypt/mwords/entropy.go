package mwords

import (
	"errors"

	"cpl.li/go/cryptor/crypt"
)

const (
	entropyMaxBits  = 256
	entropyMinBits  = 128
	entropyMultiple = 32
)

// check that the number of bits is within max, min and a multiple of 32
func isValidEntropy(bits uint) bool {
	if (bits%entropyMultiple) != 0 ||
		bits < entropyMinBits || bits > entropyMaxBits {

		return false
	}
	return true
}

// return a valid number of entropy/random bits as a byte array
func entropy(bits uint) ([]byte, error) {
	// validate input
	if !isValidEntropy(bits) {
		return nil, errors.New("invalid bit count")
	}

	// return random data
	return crypt.RandomBytes(bits / 8), nil
}
