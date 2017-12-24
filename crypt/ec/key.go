package ec

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

// PublicKey is a substitute for ecdsa.PublicKey.
type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}

// ImportPublic a Go ECDSA Public key as a custom Public key.
func ImportPublic(pub *ecdsa.PublicKey) *PublicKey {
	return &PublicKey{
		Curve: pub.Curve,
		X:     pub.X,
		Y:     pub.Y,
	}
}

// Export converts a custom ECDSA Public key to Go standard ECDSA Public key.
func (pub *PublicKey) Export() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{Curve: pub.Curve, X: pub.X, Y: pub.Y}
}

// PrivateKey is a substitute for ecdsa.PrivateKey.
type PrivateKey struct {
	PublicKey
	D *big.Int
}

// Import a Go ECDSA Private key as custom Private key.
func Import(prv *ecdsa.PrivateKey) *PrivateKey {
	pub := ImportPublic(&prv.PublicKey)
	return &PrivateKey{*pub, prv.D}
}

// Export converts a custom ECDSA Private key to Go standard ECDSA Private key.
func (prv *PrivateKey) Export() *ecdsa.PrivateKey {
	customPublicKey := &prv.PublicKey
	ecdsaPublicKey := customPublicKey.Export()
	return &ecdsa.PrivateKey{PublicKey: *ecdsaPublicKey, D: prv.D}
}

func comparePrivate(prv0, prv1 *PrivateKey) bool {
	if prv0.D.Cmp(prv1.D) != 0 {
		return false
	}

	return comparePublic(&prv0.PublicKey, &prv1.PublicKey)
}

func comparePublic(pub0, pub1 *PublicKey) bool {
	if pub0.X.Cmp(pub1.X) != 0 || pub0.Y.Cmp(pub1.Y) != 0 {
		return false
	}

	return true
}

// IsEqual takes a private key and compares it to the current private key.
func (prv *PrivateKey) IsEqual(o *PrivateKey) bool {
	return comparePrivate(prv, o)
}

// IsEqual takes a public key and compares it to the current public key.
func (pub *PublicKey) IsEqual(o *PublicKey) bool {
	return comparePublic(pub, o)
}
