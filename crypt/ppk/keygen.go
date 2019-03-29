package ppk

import (
	"crypto/rand"

	"golang.org/x/crypto/curve25519"
)

// NewPrivateKey generates a random  curve25519 secret key.
func NewPrivateKey() (sk PrivateKey, err error) {
	_, err = rand.Read(sk[:])
	return
}

// PublicKey generates the public key from the secret key.
func (sk *PrivateKey) PublicKey() (pk PublicKey) {
	curve25519.ScalarBaseMult(
		(*[KeySize]byte)(&pk),
		(*[KeySize]byte)(sk))
	return
}

// SharedSecret generates a shared secret between exchanged public keys.
func (sk *PrivateKey) SharedSecret(pk PublicKey) (ss [KeySize]byte) {
	curve25519.ScalarMult(&ss,
		(*[KeySize]byte)(sk),
		(*[KeySize]byte)(&pk))
	return ss
}
