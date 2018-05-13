package ppk_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/ppk"
)

func TestPrivateKeyEncoding(t *testing.T) {
	t.Parallel()

	key := ppk.NewKey()
	eKey := ppk.Encode(key)
	dKey := ppk.DecodePrivate(eKey)
	if !bytes.Equal(ppk.Encode(key), ppk.Encode(dKey)) {
		t.Errorf("ppk failed to encode then decode, not matching")
	}
}

func TestPublicKeyEncoding(t *testing.T) {
	t.Parallel()

	key := ppk.NewKey()
	eKey := ppk.Encode(&key.PublicKey)
	dKey := ppk.DecodePublic(eKey)
	if !bytes.Equal(ppk.Encode(&key.PublicKey), ppk.Encode(dKey)) {
		t.Errorf("ppk failed to encode then decode, not matching")
	}
}

func TestKeyIntegrity(t *testing.T) {
	t.Parallel()

	// Generate set of keys
	prv := ppk.NewKey()
	pub := &prv.PublicKey

	// Encode keys
	ePrv := ppk.Encode(prv)
	ePub := ppk.Encode(pub)

	// Decode keys
	dPrv := ppk.DecodePrivate(ePrv)
	dPub := ppk.DecodePublic(ePub)

	// Generate random data
	data := crypt.RandomData(10)

	// Encrypt data with both keys
	eData0, err := ppk.Encrypt(pub, data)
	if err != nil {
		t.Error(err)
	}
	eData1, err := ppk.Encrypt(dPub, data)
	if err != nil {
		t.Error(err)
	}

	// Decrypt data with opposite keys
	dData0, err := ppk.Decrypt(dPrv, eData0)
	if err != nil {
		t.Error(err)
	}
	dData1, err := ppk.Decrypt(prv, eData1)
	if err != nil {
		t.Error(err)
	}

	// Check data
	if !bytes.Equal(dData0, dData1) {
		t.Errorf("ppk data does not match after decryption")
	}

	// Generate signatures
	sign0, err := ppk.Sign(prv, data)
	if err != nil {
		t.Error(err)
	}
	sign1, err := ppk.Sign(dPrv, data)
	if err != nil {
		t.Error(err)
	}

	// Check signatures
	if !ppk.Verify(dPub, data, sign0) || !ppk.Verify(pub, data, sign1) {
		t.Errorf("ppk signature failed to verify")
	}
}
