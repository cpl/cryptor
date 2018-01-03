package ec

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/thee-engineer/cryptor/common/math"
	"github.com/thee-engineer/cryptor/crypt"
)

// Encode returns the a []byte representation of the key
func (prv *PrivateKey) Encode() []byte {
	return math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
}

// EncodeString returns the a hex string representation of the key
func (prv *PrivateKey) EncodeString() string {
	return crypt.EncodeString(
		math.PaddedBigBytes(prv.D, prv.Params().BitSize/8))
}

// Decode returns a key from a given []byte representation of the key
func Decode(keyData []byte) (*PrivateKey, error) {
	// Create new key
	prv := new(PrivateKey)
	prv.PublicKey.Curve = ellipticCurveFunc

	// Validate data lenght
	if 8*len(keyData) != prv.Params().BitSize {
		return nil, fmt.Errorf(
			"invalid length, need %d bits", prv.Params().BitSize)
	}

	// Set private/public keys
	prv.D = new(big.Int).SetBytes(keyData)
	prv.PublicKey.X, prv.PublicKey.Y = prv.PublicKey.Curve.ScalarBaseMult(keyData)

	// TODO: Possibly? test this
	if prv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}

	return prv, nil
}

// DecodeString takes a hex string representation of a key and returns the
// private key.
func DecodeString(keyHex string) (*PrivateKey, error) {
	keyData, err := crypt.DecodeString(keyHex)
	if err != nil {
		return nil, err
	}

	return Decode(keyData)
}
