package crypt_test

import (
	"testing"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/ppk"
)

const password = "testing"

func TestPBKDF2(t *testing.T) {
	expected :=
		"e9de130eaf35bb44fb45197f40363025db78b7198091d5b7ae3aa70fad95a140"
	expectedPub :=
		"6b12ef3a6de11fd8b160acd260a26be9b6f42be405fae2f7afe774048db96a7a"

	// derive key
	var key ppk.PrivateKey
	key = crypt.Key([]byte(password))

	// check len
	if len(key) != ppk.KeySize {
		t.Fatal("generated key is of invalid length")
	}

	// check expected key
	if key.ToHex() != expected {
		t.Fatal("derived key does not match expected key")
	}

	// check expected public key
	if key.PublicKey().ToHex() != expectedPub {
		t.Fatal("derived key public key does not match expected public key")
	}
}

func BenchmarkPBKDF2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		crypt.Key([]byte(password))
	}
}
