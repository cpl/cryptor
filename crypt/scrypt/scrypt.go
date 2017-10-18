package scrypt

import (
	"github.com/thee-engineer/cryptor/crypt"
	"golang.org/x/crypto/scrypt"
)

const keySize = 32

// Scrypt takes a password and salt and derives a 32 byte key.
func Scrypt(password string, salt []byte) []byte {
	// Derive 32 byte key.
	key, err := scrypt.Key([]byte(password), salt, 65536, 8, 1, keySize)
	if err != nil {
		panic(err)
	}

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
// 32byte key using scrypt.
func AllRandom() ([]byte, []byte, []byte) {
	// Generate random salt.
	salt := crypt.RandomData(16)

	// Generate random password.
	pass := crypt.RandomData(64)

	return Scrypt(string(pass), salt), pass, salt
}
