package hkdf

import (
	"crypto/hmac"

	"golang.org/x/crypto/blake2s"

	"cpl.li/go/cryptor/internal/crypt"
	"cpl.li/go/cryptor/internal/crypt/hashing"
)

// HMAC (hash-based message authentication code).
// https://tools.ietf.org/html/rfc2104
func HMAC(sum *[blake2s.Size]byte, key []byte, data ...[]byte) {
	mac := hmac.New(hashing.HashFunction, key)

	for _, set := range data {
		mac.Write(set)
	}

	mac.Sum(sum[:0])
}

// HKDF (hash key derivation function).
// https://tools.ietf.org/html/rfc5869
func HKDF(key, data []byte, outkeys ...*[blake2s.Size]byte) {
	var localKey [blake2s.Size]byte
	HMAC(&localKey, key, data)
	defer crypt.ZeroBytes(localKey[:])

	for index, outkey := range outkeys {
		iter := []byte{byte(index + 1)}

		if index <= 0 {
			HMAC(outkey, localKey[:], iter)
		} else {
			HMAC(outkey, localKey[:], outkeys[index-1][:], iter)
		}
	}
}
