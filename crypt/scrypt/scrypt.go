package scrypt

import (
	"log"

	"github.com/thee-engineer/cryptor/crypt"
	"golang.org/x/crypto/scrypt"
)

const n = 65536 // Cost Factor (2^16)
const r = 8     // Block size
const p = 1     // Parallelization Factor

const saltSize = 16
const randomPassSize = 128

// Scrypt takes a password and salt and derives a key.
func Scrypt(password string, salt []byte) []byte {
	// Derive 32 byte key.
	key, err := scrypt.Key([]byte(password), salt, n, r, p, crypt.KeySize)
	if err != nil {
		log.Panic(err)
	}

	return key
}

// RandomSalt applies the scrypt derivation function using a given password
// and a random 16 byte salt which is returned alongside the key..
func RandomSalt(password string) ([]byte, []byte) {
	// Generate random salt.
	salt := crypt.RandomData(saltSize)

	return Scrypt(password, salt), salt
}

// AllRandom generates both a random password and salt, then derives a new
// aes key size byte key using scrypt.
func AllRandom() (key []byte, pass []byte, salt []byte) {
	salt = crypt.RandomData(saltSize)       // Generate random salt
	pass = crypt.RandomData(randomPassSize) // Generate random password
	key = Scrypt(string(pass), salt)        // Derive key

	return key, pass, salt

}
