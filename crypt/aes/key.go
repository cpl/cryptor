package aes

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"strings"

	"github.com/thee-engineer/cryptor/crypt"
	"golang.org/x/crypto/pbkdf2"
)

// KeySize used by AES256.
const KeySize = 32

// Key is a byte array of AES256 Key size.
type Key [KeySize]byte

// NullKey key containing 32 of byte 0.
var NullKey Key = [KeySize]byte{}

// NewKey returns a new random AES256 Key.
func NewKey() (key Key) {
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return key
}

// NewKeyFromString takes a hex encoded string and returns an AES256 Key.
func NewKeyFromString(hex string) (key Key, err error) {
	// If empty string is given as key, return null key
	hex = strings.TrimSpace(hex)
	if hex == "" {
		return NullKey, nil
	}
	if len(hex) != 64 {
		return NullKey, errors.New("aes key error: invalid hex string size")
	}

	// Decode hex string and convert to bytes
	keyData, err := crypt.Decode([]byte(hex))
	if err != nil {
		return NullKey, err
	}

	copy(key[:], keyData)
	return key, nil
}

// NewKeyFromBytes takes a byte array and builds an AES256 Key
func NewKeyFromBytes(data []byte) (key Key) {
	copy(key[:], data)
	return key
}

// NewKeyFromPassword returns a valid key derived from a password string
// It uses SHA256 and iterates 100000 times. No salt is used.
func NewKeyFromPassword(password string) Key {
	return NewKeyFromBytes(
		pbkdf2.Key([]byte(password), nil, 100000, KeySize, sha256.New))
}

// Bytes returns the key as []byte array.
func (key Key) Bytes() []byte {
	return key[:]
}

// Encode returns a hex encoded []byte array.
func (key Key) Encode() []byte {
	return crypt.Encode(key.Bytes())
}

// String returns a hex encoded string
func (key Key) String() string {
	return string(key.Encode())
}
