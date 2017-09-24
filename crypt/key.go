package crypt

import (
	"crypto/rand"
	"io"
)

// AESKeySize ...
const AESKeySize = 32

// AESKey ...
type AESKey [AESKeySize]byte

// NullKey ...
var NullKey AESKey = [AESKeySize]byte{}

// NewKey ...
func NewKey() AESKey {
	key := [AESKeySize]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return key
}

// Bytes ...
func (key *AESKey) Bytes() []byte {
	return key[:]
}

// Encode ...
func (key *AESKey) Encode() []byte {
	return Encode(key.Bytes())
}

// String ...
func (key *AESKey) String() string {
	return string(key.Encode())
}
