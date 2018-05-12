package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// Sign ...
func Sign(key *rsa.PrivateKey, msg []byte) ([]byte, error) {
	return rsa.SignPSS(
		rand.Reader, key, crypto.BLAKE2b_256, hashing.Hash(msg), nil)
}

// Verify ...
func Verify(key *rsa.PublicKey, msg, sig []byte) bool {
	return nil == rsa.VerifyPSS(
		key, crypto.BLAKE2b_256, hashing.Hash(msg), sig, nil)
}

// Encrypt ...
func Encrypt(key *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(hashing.HashFunction(), rand.Reader, key, msg, nil)
}

// Decrypt ...
func Decrypt(key *rsa.PrivateKey, msg []byte) ([]byte, error) {
	return rsa.DecryptOAEP(hashing.HashFunction(), rand.Reader, key, msg, nil)
}
