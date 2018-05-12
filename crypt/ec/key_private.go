package ec

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/thee-engineer/cryptor/common/math"
	"github.com/thee-engineer/cryptor/crypt"
)

// PrivateKey is a substitute for ecdsa.PrivateKey.
type PrivateKey struct {
	PublicKey
	D *big.Int
}

// Import a Go ECDSA Private key as custom Private key.
func (prv *PrivateKey) Import(ecdsaKey *ecdsa.PrivateKey) {
	prv.D = ecdsaKey.D
	prv.PublicKey.Import(&ecdsaKey.PublicKey)
}

// Export converts a custom ECDSA Private key to Go standard ECDSA Private key.
func (prv *PrivateKey) Export() *ecdsa.PrivateKey {
	customPublicKey := &prv.PublicKey
	ecdsaPublicKey := customPublicKey.Export()
	return &ecdsa.PrivateKey{PublicKey: *ecdsaPublicKey, D: prv.D}
}

// Encode returns the a []byte representation of the key
func (prv *PrivateKey) Encode() []byte {
	return math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
}

// EncodeString returns the a hex string representation of the key
func (prv *PrivateKey) EncodeString() string {
	return crypt.EncodeString(prv.Encode())
}

// Decode returns a key from a given []byte representation of the key
func (prv *PrivateKey) Decode(keyData []byte) error {
	prv.PublicKey.Curve = ellipticCurveFunc

	// Validate data length
	if 8*len(keyData) != prv.Params().BitSize {
		return fmt.Errorf(
			"invalid length, need %d bits", prv.Params().BitSize)
	}

	// Set private/public keys
	prv.D = new(big.Int).SetBytes(keyData)
	prv.PublicKey.X, prv.PublicKey.Y = prv.PublicKey.Curve.ScalarBaseMult(keyData)

	if prv.PublicKey.X == nil {
		return errors.New("invalid private key")
	}

	return nil
}

// DecodeString takes a hex string representation of a key and returns the
// private key.
func (prv *PrivateKey) DecodeString(keyHex string) error {
	keyBytes, err := crypt.DecodeString(keyHex)
	if err != nil {
		return err
	}

	return prv.Decode(keyBytes)
}

// IsEqual takes a private key and compares it to the current private key.
func (prv *PrivateKey) IsEqual(o *PrivateKey) bool {
	return comparePrivate(prv, o)
}

func comparePrivate(prv0, prv1 *PrivateKey) bool {
	if prv0.D.Cmp(prv1.D) != 0 {
		return false
	}

	return comparePublic(&prv0.PublicKey, &prv1.PublicKey)
}
