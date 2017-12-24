package scrypt

import (
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"golang.org/x/crypto/scrypt"
)

const n = 65536
const r = 8
const p = 1

// Scrypt takes a password and salt and derives a aes keysize byte key.
func Scrypt(password string, salt []byte) []byte {
	// Derive 32 byte key.
	key, _ := scrypt.Key([]byte(password), salt, n, r, p, aes.KeySize)

	return key
}

// RandomSalt applies the scrypt derivation function using a given password
// and a random 16 byte salt.
func RandomSalt(password string) ([]byte, []byte) {
	// Generate random salt.
	salt := crypt.RandomData(16)

	return Scrypt(password, salt), salt
}

// AllRandom generates both a random password and salt, then derives a new
// aes key size byte key using scrypt.
func AllRandom() ([]byte, []byte, []byte) {
	// Generate random salt.
	salt := crypt.RandomData(16)

	// Generate random password.
	pass := crypt.RandomData(64)

	return Scrypt(string(pass), salt), pass, salt
}
