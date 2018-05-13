package ppk_test

import (
	"bytes"
	"testing"

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
