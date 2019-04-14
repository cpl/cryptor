package crypt_test

import (
	"testing"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/ppk"
)

const password = "testing"
const salt = ".-_cryptor,$"

func TestPBKDF2(t *testing.T) {
	expected :=
		"28df0b93627d5b50ed4fef574e774a00ac634cbd3395d0a57e769581e806f82f"
	expectedPub :=
		"a5f686a01f0327c2a1bce2d2ae01c4174d1637fd31a5a065d0b235ea37cc3d74"

	// derive key
	var key ppk.PrivateKey
	key = crypt.Key([]byte(password), []byte(salt))

	// check len
	if len(key) != ppk.KeySize {
		t.Fatal("generated key is of invalid length")
	}

	// check expected key
	if key.ToHex() != expected {
		t.Error("got", key.ToHex())
		t.Error("expected", expected)
		t.Fatal("derived key does not match expected key")
	}

	// check expected public key
	if key.PublicKey().ToHex() != expectedPub {
		t.Error("got", key.PublicKey().ToHex())
		t.Error("expected", expectedPub)
		t.Fatal("derived key public key does not match expected public key")
	}

	// check default salt is working
	dKey := crypt.Key([]byte(password), nil)
	if !key.Equals(dKey) {
		t.Fatal("default salt failed, keys don't match")
	}
}

func BenchmarkPBKDF2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		crypt.Key([]byte(password), []byte(salt))
	}
}
