package ec_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestECDSASign(t *testing.T) {
	t.Parallel()

	msg := []byte("hello world")

	// Generate two key paris
	key0, key1, err := generateKeyParis()
	if err != nil {
		t.Error(err)
	}

	// Generate first signature
	sign0, err := key0.Sign(msg)
	if err != nil {
		t.Error(err)
	}

	// Generate second signature
	sign1, err := key1.Sign(msg)
	if err != nil {
		t.Error(err)
	}

	// Compare the two signatures
	if bytes.Equal(sign0, sign1) {
		t.Error("ecdsa: unexpected signature match")
	}

	// Verify valid signatures
	if !key0.PublicKey.Verify(msg, sign0) {
		t.Error("ecdsa: failed to verify valid signature")
	}
	if !key1.PublicKey.Verify(msg, sign1) {
		t.Error("ecdsa: failed to verify valid signature")
	}

	// Verify invalid signs
	if key0.PublicKey.Verify(msg, sign1) {
		t.Error("ecdsa: unexpected signature validation")
	}
	if key1.PublicKey.Verify(msg, sign0) {
		t.Error("ecdsa: unexpected signature validation")
	}

	// Verify invalid data
	if key0.PublicKey.Verify(crypt.RandomData(100), sign1) {
		t.Error("ecdsa: unexpected data validation")
	}
	if key1.PublicKey.Verify(crypt.RandomData(30), sign0) {
		t.Error("ecdsa: unexpected data validation")
	}

	// Verify random sign
	if key0.PublicKey.Verify(msg, crypt.RandomData(64)) {
		t.Error("ecdsa: unexpected data validation")
	}
	if key1.PublicKey.Verify(msg, crypt.RandomData(100)) {
		t.Error("ecdsa: unexpected data validation")
	}
	if key0.PublicKey.Verify(msg, crypt.RandomData(2)) {
		t.Error("ecdsa: unexpected data validation")
	}
}
