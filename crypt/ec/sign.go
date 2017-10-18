package ec

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"

	"github.com/thee-engineer/cryptor/crypt"

	"github.com/thee-engineer/cryptor/crypt/hashing"
)

const signatureSize = 64

// Sign takes a message, generates the SHA256 hash and signs it using the
// ecdsa private key.
func (prv *PrivateKey) Sign(msg []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, prv.Export(), hashing.SHA256Digest(msg))
	if err != nil {
		return nil, err
	}

	// Transform big ints into bytes, and prepare signature
	rBytes, sBytes := r.Bytes(), s.Bytes()

	// Prep memory zeroing
	defer crypt.ZeroBytes(rBytes, sBytes)
	r = big.NewInt(0)
	s = big.NewInt(0)

	// Allocate byte array for signature
	signature := make([]byte, signatureSize)

	// Move R, S bytes into signature
	copy(signature[:signatureSize/2], rBytes)
	copy(signature[signatureSize/2:], sBytes)

	return signature, nil
}

// Verify takes a message and the supposed signature, it generates the
// hash of the msg and tries validates the msg using the ecdsa public key.
func (pub *PublicKey) Verify(msg, signature []byte) bool {
	// Check for valid signature
	if len(signature) != signatureSize {
		return false
	}

	// Convert bytes back into big ints
	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:signatureSize/2])
	s.SetBytes(signature[signatureSize/2:])

	// Use Go's ecdsa to verify signature
	return ecdsa.Verify(pub.Export(), hashing.SHA256Digest(msg), r, s)
}
