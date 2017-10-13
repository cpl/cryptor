package crypt

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"hash"
	"math/big"

	"golang.org/x/crypto/sha3"
)

// PublicKey is a substitute for ecdsa.PublicKey.
type PublicKey struct {
	elliptic.Curve
	X, Y   *big.Int
	Params *Param
}

// ImportPublic a Go ECDSA Public key as a custom Public key.
func ImportPublic(pub *ecdsa.PublicKey) *PublicKey {
	return &PublicKey{
		Curve:  pub.Curve,
		X:      pub.X,
		Y:      pub.Y,
		Params: ECDSA521Param,
	}
}

// Export converts a custom ECDSA Public key to Go standard ECDSA Public key.
func (pub *PublicKey) Export() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{Curve: pub.Curve, X: pub.X, Y: pub.Y}
}

// Param is information about the ECDSA Curve
type Param struct {
	Hash      func() hash.Hash
	hType     crypto.Hash
	Cipher    func([]byte) (cipher.Block, error)
	BlockSize int
	KeyLength int
}

var (
	// ECDSA521Param are all the standard parameters for Curve 521.
	ECDSA521Param = &Param{
		Hash:      sha3.New512,   // Using 512 sha3 function
		hType:     crypto.SHA512, // Using 512 sha3
		Cipher:    aes.NewCipher, // Asymetric cipher block
		BlockSize: aes.BlockSize, // Asymetric cipher block size
		KeyLength: 32}            // Length of the asymetric key

	// ECDSA256Param are all the standard parameters for Curve 251.
	ECDSA256Param = &Param{
		Hash:      sha256.New,
		hType:     crypto.SHA256,
		Cipher:    aes.NewCipher,
		BlockSize: aes.BlockSize,
		KeyLength: 32,
	}
)

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

func generateKey(curve elliptic.Curve) (*PrivateKey, error) {
	pb, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	prv := new(PrivateKey)
	prv.PublicKey.X = x
	prv.PublicKey.Y = y
	prv.PublicKey.Curve = curve
	prv.D = new(big.Int).SetBytes(pb)

	return prv, nil
}

// GenerateKey256 generates a new Private/Public key pair using P256 Curve
// which implements P-256 (see).
func GenerateKey256() (*PrivateKey, error) {
	return generateKey(elliptic.P256())
}

// GenerateKey521 generates a new Private/Public key pair using P521 Curve
// which implements P-521 (see FIPS 186-3, section D.2.5).
func GenerateKey521() (*PrivateKey, error) {
	return generateKey(elliptic.P521())
}

// GenerateShared creates a new shared secret of given lenght.
func (prv *PrivateKey) GenerateShared(pub *PublicKey, secLen, macLen int) (sec []byte, err error) {
	// Verify matching curves on both public and private keys
	if prv.PublicKey.Curve != pub.Curve {
		return nil, errors.New("shared key: invalid curve match")
	}

	// Verify shared secret length
	if secLen+macLen > pub.SecretLength() {
		return nil, errors.New("shared key: length too big")
	}

	// Verify new point position
	np, _ := pub.Curve.ScalarMult(pub.X, pub.Y, prv.D.Bytes())
	if np == nil {
		return nil, errors.New("shared key: invalid point position")
	}

	// Generate shared secret
	sec = make([]byte, secLen+macLen)
	secBytes := np.Bytes()
	copy(sec[len(sec)-len(secBytes):], secBytes)
	return sec, nil
}

// SecretLength returns the maximum length of the shared secret.
func (pub *PublicKey) SecretLength() int {
	return (pub.Curve.Params().BitSize + 7) / 8
}
