package scrypt

import (
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/utils"
	"golang.org/x/crypto/scrypt"
)

const n = 65536 // Cost Factor (2^16)
const r = 8     // Block size
const p = 1     // Parallelization Factor

// SaltSize ...
const SaltSize = 16
const randomPassSize = 128

// Scrypt takes a password and salt and derives a key.
func Scrypt(password string, salt []byte) []byte {
	// Derive 32 byte key.
	key, err := scrypt.Key([]byte(password), salt, n, r, p, crypt.KeySize)
	utils.CheckErr(err)
	return key
}

// RandomSalt applies the scrypt derivation function using a given password
// and a random 16 byte salt which is returned alongside the key..
func RandomSalt(password string) ([]byte, []byte) {
	// Generate random salt.
	salt := crypt.RandomData(SaltSize)

	return Scrypt(password, salt), salt
}

// AllRandom generates both a random password and salt, then derives a new
// aes key size byte key using scrypt.
func AllRandom() (key []byte, pass []byte, salt []byte) {
	salt = crypt.RandomData(SaltSize)       // Generate random salt
	pass = crypt.RandomData(randomPassSize) // Generate random password
	key = Scrypt(string(pass), salt)        // Derive key

	return key, pass, salt

}
