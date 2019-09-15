package hkdf

import (
	"crypto/hmac"

	"golang.org/x/crypto/blake2s"

	"cpl.li/go/cryptor/crypt"
)

// HMAC generates the "key-hash message authentication code" using BLAKE2s hash.
func HMAC(sum *[blake2s.Size]byte, key []byte, data ...[]byte) {
	// generate MAC
	mac := hmac.New(crypt.HashFunction, key)

	// iterate byte arrays
	for _, set := range data {
		mac.Write(set)
	}

	// put sum at pointer
	mac.Sum(sum[:0])
}

// HKDF applies the HMAC algorithm with blake2s hashing and generates `count`
// number of keys returned as a list of byte arrays.
func HKDF(key, data []byte, outkeys ...*[blake2s.Size]byte) {
	// generate key used in all future HMAC instances
	var localKey [blake2s.Size]byte
	HMAC(&localKey, key, data)
	defer crypt.ZeroBytes(localKey[:])

	// iterate and generate new keys
	for index, outkey := range outkeys {
		// iter will be a byte in the range [0x1, 0x2, 0x3 ... 0xCOUNT]
		iter := []byte{byte(index + 1)}

		// if the first generated key, there is no prev to use
		if index <= 0 {
			HMAC(outkey, localKey[:], iter)
		} else {
			HMAC(outkey, localKey[:], outkeys[index-1][:], iter)
		}
	}
}
