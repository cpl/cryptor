package scrypt

import (
	"github.com/thee-engineer/cryptor/crypt"
	"golang.org/x/crypto/scrypt"
)

const n = 65536 // Cost Factor (2^16)
const r = 8     // Block size
const p = 1     // Parallelization Factor

const saltSize = 16
const passSize = 128

// Scrypt takes a password and salt and derives a aes keysize byte key.
func Scrypt(password string, salt []byte) []byte {
	// Derive 32 byte key.
	key, err := scrypt.Key([]byte(password), salt, n, r, p, crypt.KeySize)
	if err != nil {
		panic(err)
	}

	return key
}

// RandomSalt applies the scrypt derivation function using a given password
// and a random 16 byte salt.
func RandomSalt(password string) ([]byte, []byte) {
	// Generate random salt.
	salt := crypt.RandomData(saltSize)

	return Scrypt(password, salt), salt
}

// AllRandom generates both a random password and salt, then derives a new
// aes key size byte key using scrypt.
func AllRandom() ([]byte, []byte, []byte) {
	// Generate random salt.
	salt := crypt.RandomData(saltSize)

	// Generate random password.
	pass := crypt.RandomData(passSize)

	return Scrypt(string(pass), salt), pass, salt
}
