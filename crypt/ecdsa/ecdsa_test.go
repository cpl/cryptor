package crypt_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func compareKeys(key0, key1 *crypt.PrivateKey) bool {
	if !bytes.Equal(key0.D.Bytes(), key1.D.Bytes()) {
		return false
	}
	if !bytes.Equal(key0.PublicKey.X.Bytes(), key1.PublicKey.X.Bytes()) {
		return false
	}
	if !bytes.Equal(key0.PublicKey.Y.Bytes(), key1.PublicKey.Y.Bytes()) {
		return false
	}

	return true
}

func generateKeyPair256() (prv0, prv1 *crypt.PrivateKey, err error) {
	prv0, err = crypt.GenerateKey256()
	if err != nil {
		return nil, nil, err
	}

	prv1, err = crypt.GenerateKey256()
	if err != nil {
		return nil, nil, err
	}

	return prv0, prv1, nil
}

func TestImportExport(t *testing.T) {
	t.Parallel()

	// Generate key pair
	prv, err := crypt.GenerateKey256()
	if err != nil {
		t.Error(err)
	}
	pub := prv.PublicKey

	// Export to Go ECDSA Keys
	ecdsaPrv := prv.Export()
	ecdsaPub := pub.Export()

	// Import private key
	impPrv := crypt.Import(ecdsaPrv)

	if !bytes.Equal(impPrv.D.Bytes(), prv.D.Bytes()) {
		t.Error("ecdsa: imported and inital key mismatch, private")
	}

	// Import public key
	impPub := crypt.ImportPublic(ecdsaPub)

	if !bytes.Equal(impPub.X.Bytes(), pub.X.Bytes()) {
		t.Error("ecdsa: imported and inital key mismatch, public")
	}

	if !bytes.Equal(impPub.Y.Bytes(), pub.Y.Bytes()) {
		t.Error("ecdsa: imported and inital key mismatch, public")
	}
}

func TestSharedSecret(t *testing.T) {
	t.Parallel()

	prv0, prv1, err := generateKeyPair256()
	if err != nil {
		t.Error(err)
	}
	pub0 := &prv0.PublicKey
	pub1 := &prv1.PublicKey

	// Obtain shared secret from first private key and second public key
	ss0, err := prv0.GenerateShared(pub1, 16, 16)
	if err != nil {
		t.Error(err)
	}

	// Obtain shared secret from second private key and first public key
	ss1, err := prv1.GenerateShared(pub0, 16, 16)
	if err != nil {
		t.Error(err)
	}

	// Compare the two secrets
	if !bytes.Equal(ss0, ss1) {
		t.Errorf("shared secret: not matching")
	}
}

func Test521Key(t *testing.T) {
	t.Parallel()

	// Generate 521 curve keys
	key0, err := crypt.GenerateKey521()
	if err != nil {
		t.Error(err)
	}

	key1, err := crypt.GenerateKey521()
	if err != nil {
		t.Error(err)
	}

	// Generate secrets on each "side"
	ss0, err := key0.GenerateShared(&key1.PublicKey, 33, 33)
	if err != nil {
		t.Error(err)
	}

	ss1, err := key1.GenerateShared(&key0.PublicKey, 33, 33)
	if err != nil {
		t.Error(err)
	}

	// Compare the two secrets
	if !bytes.Equal(ss0, ss1) {
		t.Errorf("shared secret: not matching")
	}
}

func TestSecretToKey(t *testing.T) {
	t.Parallel()

	prv0, prv1, err := generateKeyPair256()
	if err != nil {
		t.Error(err)
	}
	pub0 := &prv0.PublicKey
	pub1 := &prv1.PublicKey

	// Obtain shared secret from first private key and second public key
	ss0, err := prv0.GenerateShared(pub1, 16, 16)
	if err != nil {
		t.Error(err)
	}

	// Obtain shared secret from second private key and first public key
	ss1, err := prv1.GenerateShared(pub0, 16, 16)
	if err != nil {
		t.Error(err)
	}

	// Compare the two secrets
	if !bytes.Equal(ss0, ss1) {
		t.Errorf("shared secret: not matching")
	}

	ss0Key, err := crypt.Key256FromSecret(ss0)
	if err != nil {
		t.Error(err)
	}

	ss1Key, err := crypt.Key256FromSecret(ss1)
	if err != nil {
		t.Error(err)
	}

	if !compareKeys(ss0Key, ss1Key) {
		t.Error("shared keys: not matching")
	}
}

func TestSharedSecretErrors(t *testing.T) {
	t.Parallel()

	prv0, prv1, err := generateKeyPair256()
	if err != nil {
		t.Error(err)
	}
	pub1 := &prv1.PublicKey

	key521, err := crypt.GenerateKey521()
	if err != nil {
		t.Error(err)
	}

	_, err = prv0.GenerateShared(pub1, 100, 0)
	if err == nil {
		t.Error("shared keys: generating invalid size key")
	}

	_, err = crypt.Key256FromSecret(crypt.RandomData(100))
	if err == nil {
		t.Error("shared keys: generating invalid size key")
	}

	_, err = prv0.GenerateShared(&key521.PublicKey, 33, 33)
	if err == nil {
		t.Error("shared keys: invalid public key type")
	}
}
