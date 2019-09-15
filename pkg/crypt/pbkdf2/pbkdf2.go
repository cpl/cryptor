package pbkdf2

import (
	"cpl.li/go/cryptor/pkg/crypt/hashing"
	"cpl.li/go/cryptor/pkg/crypt/ppk"
	"golang.org/x/crypto/pbkdf2"
)

const staticSalt = ".-_cryptor,$"
const iter = 131072

// Key takes a password + salt and applies pbkdf2 in order to derive a key. This
// uses a custom number of iterations, hash function and key size. PBKDF2 is
// defined in RFC 2898.
func Key(password, salt []byte) (key [ppk.KeySize]byte) {
	// default salt
	if salt == nil {
		salt = []byte(staticSalt)
	}

	copy(
		key[:],
		pbkdf2.Key(password, salt, iter, ppk.KeySize, hashing.HashFunction))
	return
}
