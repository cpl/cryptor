package crypt_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestImportExport(t *testing.T) {
	t.Parallel()

	// Generate key pair
	prv, err := crypt.GenerateKey()
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

	// Generate first key pair
	prv0, err := crypt.GenerateKey()
	if err != nil {
		t.Error(err)
	}
	pub0 := &prv0.PublicKey

	// Generate second key pair
	prv1, err := crypt.GenerateKey()
	if err != nil {
		t.Error(err)
	}
	pub1 := &prv1.PublicKey

	// Obtain shared secret from first private key and second public key
	ss0, err := prv0.GenerateShared(pub1)
	if err != nil {
		t.Error(err)
	}

	// Obtain shared secret from second private key and first public key
	ss1, err := prv1.GenerateShared(pub0)
	if err != nil {
		t.Error(err)
	}

	// Compare the two secrets
	if !bytes.Equal(ss0, ss1) {
		t.Errorf("shared secret: not matching")
	}
}
