package ec_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"

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
	prv := new(ec.PrivateKey)
	prv.Import(ecdsaKey)
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
	key0Clone := new(ec.PrivateKey)
	key0Clone.Import(key0.Export())

	// Check for equal keys on different keys
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

	// Compare two different public keys
	if key0.PublicKey.IsEqual(&key1.PublicKey) {
		t.Error("ecdsa: unexpected public key equality")
	}
}

func TestKeyEncoding(t *testing.T) {
	t.Parallel()

	key, err := ec.GenerateKey()
	if err != nil {
		t.Error(err)
	}

	outByte := key.Encode()
	outString := key.EncodeString()

	if crypt.EncodeString(outByte) != outString {
		t.Log(crypt.EncodeString(outByte))
		t.Log(outString)
		t.Error("ec key | mismatch key encodings")
	}

	decodedKey := new(ec.PrivateKey)
	if err := decodedKey.Decode(outByte); err != nil {
		t.Error(err)
	}

	if !key.IsEqual(decodedKey) {
		t.Log("decoded: ", decodedKey)
		t.Log("original:", key)
		t.Errorf("ec key | mismatch original with decoded key")
	}

	if err := decodedKey.DecodeString(outString); err != nil {
		t.Error(err)
	}

	if !key.IsEqual(decodedKey) {
		t.Log("decoded: ", decodedKey)
		t.Log("original:", key)
		t.Errorf("ec key | mismatch original with decoded key")
	}

	if err := decodedKey.Decode([]byte{10, 20, 30}); err == nil {
		t.Error("ec key | decoded invalid []byte")
	}

	if err := decodedKey.Decode(crypt.RandomData(65)); err == nil {
		t.Error("ec key | decoded invalid []byte")
	}

	if err := decodedKey.DecodeString("testing"); err == nil {
		t.Error("ec key | decoded invalid string")
	}
}

func TestPublicKeyEncoding(t *testing.T) {
	prv, _ := ec.GenerateKey()
	pub := &prv.PublicKey

	if pub.EncodeString() != crypt.EncodeString(pub.Encode()) {
		t.Errorf("ec key | hex encoding does not match")
	}

	npub := new(ec.PublicKey)
	if err := npub.Decode(pub.Encode()); err != nil {
		t.Error(err)
	}

	if !pub.IsEqual(npub) {
		t.Errorf("ec key | keys do not match after encoding/decoding")
	}

	if err := npub.DecodeString(pub.EncodeString()); err != nil {
		t.Error(err)
	}

	if !pub.IsEqual(npub) {
		t.Errorf("ec key | keys do not match after encoding/decoding")
	}
}
