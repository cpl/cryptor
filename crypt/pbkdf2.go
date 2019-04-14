package crypt

import (
	"cpl.li/go/cryptor/crypt/ppk"
	"golang.org/x/crypto/pbkdf2"
)

// TODO Make this work with randomly generated salts

const salt = ".-_cryptor,$"
const iter = 131072

// Key takes a password and applies pbkdf2 in order to derive a key.
// PBKDF2 is defined in RFC 2898.
func Key(password []byte) (key [ppk.KeySize]byte) {
	copy(
		key[:],
		pbkdf2.Key(password, []byte(salt), iter, ppk.KeySize, HashFunction))
	return
}
