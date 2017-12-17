// Package ec implements elliptic curve private/public key pairs using P256.
// This package also allows for the creation of shared secret keys, which can
// be used in combination with Symetric encryption ciphers.
// Relevant XKCD https://xkcd.com/538/.
package ec

import (
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/thee-engineer/cryptor/crypt"
)

const secretSize = 64

var ellipticCurveFunc = elliptic.P256()

// GenerateKey generates a new Private(+Public) ecdsa key pair using P256 Curve
// which implements P-256.
func GenerateKey() (*PrivateKey, error) {
	pb, x, y, err := elliptic.GenerateKey(ellipticCurveFunc, rand.Reader)
	if err != nil {
		return nil, err
	}
	defer crypt.ZeroBytes(pb)

	// Allocate memory for PrivateKey
	prv := new(PrivateKey)

	// Allocate public key values
	prv.PublicKey.X = x
	prv.PublicKey.Y = y

	// Assign ec function
	prv.PublicKey.Curve = ellipticCurveFunc

	// Set private key bytes
	prv.D = new(big.Int).SetBytes(pb)

	return prv, nil
}

// GenerateSecret creates a new shared secret.
func (prv *PrivateKey) GenerateSecret(pub *PublicKey) ([]byte, error) {
	// Verify matching curves on both public and private keys
	if prv.PublicKey.Curve != pub.Curve {
		return nil, errors.New("shared key: invalid curve match")
	}

	// Verify new point position
	point, _ := pub.Curve.ScalarMult(pub.X, pub.Y, prv.D.Bytes())
	if point == nil {
		return nil, errors.New("shared key: invalid point position")
	}

	// Convert elliptic curve point to bytes
	secret := point.Bytes()

	// Prepare to zero secret from memory
	point = big.NewInt(0)

	return secret, nil
}
