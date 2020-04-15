package ppk

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

// NewPrivateKey ...
func NewPrivateKey(sk *PrivateKey) error {
	if sk != nil {
		_, err := rand.Read(sk[:])
		return err
	}

	return fmt.Errorf("must provide key pointer")
}

// PublicKey ...
func (sk *PrivateKey) PublicKey(pk *PublicKey) error {
	if pk != nil {
		curve25519.ScalarBaseMult(
			(*[KeySize]byte)(pk),
			(*[KeySize]byte)(sk))
		return nil
	}

	return fmt.Errorf("must provide key pointer")
}

// SharedSecret ...
func (sk *PrivateKey) SharedSecret(pk *PublicKey, ss *[KeySize]byte) {
	curve25519.ScalarMult(ss,
		(*[KeySize]byte)(sk),
		(*[KeySize]byte)(pk))

	return
}
