package pbkdf2

import (
	"golang.org/x/crypto/pbkdf2"

	"cpl.li/go/cryptor/internal/crypt/hashing"
	"cpl.li/go/cryptor/internal/crypt/ppk"
)

const (
	staticSalt = ".-_cryptor,$"
	iter       = 131072
)

// Key will derive a 32byte key from the given password and salt using pbkdf2,
// the Cryptor hashing.HashFunction. If a nil salt is given, the default salt
// is used.
func Key(password, salt []byte) (key [ppk.KeySize]byte) {
	if salt == nil {
		salt = []byte(staticSalt)
	}

	copy(
		key[:],
		pbkdf2.Key(password, salt, iter, ppk.KeySize, hashing.HashFunction))
	return
}
