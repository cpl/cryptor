package ec

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"

	"github.com/thee-engineer/cryptor/crypt"
)

// PublicKey is a substitute for ecdsa.PublicKey.
type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}

// Import a Go ECDSA Public key as a custom Public key.
func (pub *PublicKey) Import(ecdsaKey *ecdsa.PublicKey) {
	pub.Curve = ecdsaKey.Curve
	pub.X = ecdsaKey.X
	pub.Y = ecdsaKey.Y
}

// Export converts a custom ECDSA Public key to Go standard ECDSA Public key.
func (pub *PublicKey) Export() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{Curve: pub.Curve, X: pub.X, Y: pub.Y}
}

// Encode ...
func (pub *PublicKey) Encode() []byte {
	return elliptic.Marshal(ellipticCurveFunc, pub.X, pub.Y)
}

// EncodeString ...
func (pub *PublicKey) EncodeString() string {
	return crypt.EncodeString(pub.Encode())
}

// Decode ...
func (pub *PublicKey) Decode(keyData []byte) error {
	pub.Curve = ellipticCurveFunc
	pub.X, pub.Y = elliptic.Unmarshal(ellipticCurveFunc, keyData)
	return nil
}

// DecodeString ...
func (pub *PublicKey) DecodeString(keyHex string) error {
	keyBytes, err := crypt.DecodeString(keyHex)
	if err != nil {
		return err
	}

	return pub.Decode(keyBytes)
}

// IsEqual takes a public key and compares it to the current public key.
func (pub *PublicKey) IsEqual(o *PublicKey) bool {
	return comparePublic(pub, o)
}

func comparePublic(pub0, pub1 *PublicKey) bool {
	if pub0.X.Cmp(pub1.X) != 0 || pub0.Y.Cmp(pub1.Y) != 0 {
		return false
	}

	return true
}
