package crypt

import (
	"crypto/rand"
	"io"
)

// AESKeySize used by AES256
const AESKeySize = 32

// AESKey is a byte array of AES256 Key sizes
type AESKey [AESKeySize]byte

// NullKey key containing 32 of byte 0
var NullKey AESKey = [AESKeySize]byte{}

// NewKey returns a new random AES256 Key
func NewKey() AESKey {
	key := [AESKeySize]byte{}

	// TODO: Add option for selecting rand.Reader
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return key
}

// NewKeyFromString takes a hex encoded string and returns a AES256 Key
func NewKeyFromString(hex string) (key AESKey) {
	// If empty string is given as key, return null key
	if hex == "" || hex == " " {
		return NullKey
	}

	// Decode hex string and convert to bytes
	keyData, err := Decode([]byte(hex))
	if err != nil {
		panic(err)
	}
	copy(key[:], keyData)
	return key
}

// NewKeyFromBytes takes a byte array and builds an AES256 Key
func NewKeyFromBytes(data []byte) (key AESKey) {
	copy(key[:], data)
	return key
}

// Bytes returns the key as []byte
func (key AESKey) Bytes() []byte {
	return key[:]
}

// Encode returns a hex encoded []byte
func (key AESKey) Encode() []byte {
	return Encode(key.Bytes())
}

// String returns a hex encoded string
func (key AESKey) String() string {
	return string(key.Encode())
}
