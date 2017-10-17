package ec_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/thee-engineer/cryptor/crypt/ec"
)

func TestECDSAKeysImportExport(t *testing.T) {
	t.Parallel()

	// Generate Go ECDSA key pair
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err)
	}

	// Import Go ECDSA key as custom ec key
	prv := ec.Import(ecdsaKey)
	pub := prv.PublicKey

	// Export custom ec key back to Go ECDSA
	ecdsaKeyExport := prv.Export()
	ecdsaPubKeyExport := pub.Export()

	// Compare private keys
	if ecdsaKey.D != ecdsaKeyExport.D {
		t.Error("ecdsa: exported key does not match")
	}

	// Compare public keys
	if ecdsaPubKeyExport.X != ecdsaKey.PublicKey.X &&
		ecdsaPubKeyExport.Y != ecdsaKey.PublicKey.Y {
		t.Error("ecdsa: exported key does not match")
	}
}

func TestECDSACompare(t *testing.T) {
	t.Parallel()

	// Generate key pairs
	key0, key1, err := generateKeyParis()
	if err != nil {
		t.Error(err)
	}
	// Clone one of the keys
	key0Clone := ec.Import(key0.Export())

	// Check for equal keys on diffrent keys
	if key0.IsEqual(key1) || key1.IsEqual(key0) {
		t.Error("ecdsa: unexpected key equality")
	}

	// Compare two equal keys
	if !key0.IsEqual(key0Clone) {
		t.Error("ecdsa: failed to find equal keys")
	}

	// Compare two equal public keys
	if !key0.PublicKey.IsEqual(&key0Clone.PublicKey) {
		t.Error("ecdsa: failed to find public equal keys")
	}

	// Compare two diffrent public keys
	if key0.PublicKey.IsEqual(&key1.PublicKey) {
		t.Error("ecdsa: unexpected public key equality")
	}
}
