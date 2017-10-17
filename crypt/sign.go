package crypt

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
)

// Sign signs a msg of any size by creating the SHA256 hash of the msg and
// applying ECDSA 256 curve signature.
func (p *PrivateKey) Sign(msg []byte) ([]byte, error) {
	// Hash message
	hash := SHA256Data(msg)

	// Go ECDSA key
	ecdsaKey := p.Export()

	// Sign message
	r, s, err := ecdsa.Sign(rand.Reader, ecdsaKey, hash.Sum(nil))
	if err != nil {
		return nil, err
	}

	// Get curve byte size
	curveSize := ecdsaKey.Params().BitSize / 8

	// Convert big ints to bytes
	rBytes, sBytes := r.Bytes(), s.Bytes()

	// Create byte array for signature
	signature := make([]byte, curveSize*2)

	// Move the R, S bytes into the signature
	copy(signature[curveSize-len(rBytes):], rBytes)
	copy(signature[curveSize*2-len(sBytes):], sBytes)

	return signature, nil
}

// Verify checks an ECDSA signature. Returns true if valid, false if invalid.
func (p *PublicKey) Verify(data, signature []byte) bool {
	// Hash message
	hash := SHA256Data(data).Sum(nil)

	// Get curve byte size
	curveSize := p.Export().Params().BitSize / 8

	// Convert bytes back into big ints
	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:curveSize])
	s.SetBytes(signature[curveSize:])

	// Use Go's ecdsa to verify signature
	return ecdsa.Verify(p.Export(), hash[:], r, s)
}
