package aes

import (
	"crypto/rand"
	"errors"
	"io"
	"strings"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/encode/b16"
	"github.com/thee-engineer/cryptor/crypt/scrypt"
)

// Key is a byte array of AES256 Key size.
type Key [crypt.KeySize]byte

// NullKey key containing 32 of byte 0.
var NullKey Key = [crypt.KeySize]byte{0}

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
		return NullKey, errors.New("aes key: invalid hex string, empty")
	}

	if len(hex) != 64 {
		return NullKey, errors.New("aes key error: invalid hex string size")
	}

	// Decode hex string and convert to bytes
	keyData, err := b16.Decode([]byte(hex))
	if err != nil {
		return NullKey, err
	}

	return NewKeyFromBytes(keyData)
}

// NewKeyFromBytes takes a byte array and builds an AES256 Key
func NewKeyFromBytes(data []byte) (key Key, err error) {
	if len(data) != crypt.KeySize {
		return NullKey, errors.New("aes key: invalid []byte len for new key")
	}
	copy(key[:], data)

	return key, nil
}

// NewKeyFromPassword returns a valid key derived from a password string.
// It uses scrypt with no salt.
func NewKeyFromPassword(password string) (Key, error) {
	return NewKeyFromBytes(scrypt.Scrypt(password, nil))
}

// Bytes returns the key as []byte array.
func (key Key) Bytes() []byte {
	return key[:]
}

// Encode returns a hex encoded []byte array.
func (key Key) Encode() []byte {
	return b16.Encode(key.Bytes())
}

// String returns a hex encoded string
func (key Key) String() string {
	return string(key.Encode())
}

// IsEqual compares two keys and returns true if the bytes match, false otherwise
// func (key Key) IsEqual(other Key) bool {
// 	return bytes.Equal(key.Bytes(), other.Bytes())
// }
