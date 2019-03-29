package ppk_test

import (
	"encoding/hex"
	"testing"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/tests"
)

func TestSharedSecretGeneration(t *testing.T) {
	// generate two private keys
	sk0, err := ppk.NewPrivateKey()
	tests.AssertNil(t, err)
	sk1, err := ppk.NewPrivateKey()
	tests.AssertNil(t, err)

	// generate their respective public keys
	pk0 := sk0.PublicKey()
	pk1 := sk1.PublicKey()

	// compute the shared secret for each pair
	ss0 := sk0.SharedSecret(pk1)
	ss1 := sk1.SharedSecret(pk0)

	// match the two shared secrets
	if ss0 != ss1 {
		// if failed, log all data
		t.Logf("unexpected mismatch of shared secrets, dumping data:\n")

		// private keys
		t.Logf("sk0: %s\n", sk0.ToHex())
		t.Logf("sk1: %s\n", sk1.ToHex())

		// public keys
		t.Logf("pk0: %s\n", pk0.ToHex())
		t.Logf("pk1: %s\n", pk1.ToHex())

		// shared secrets
		t.Logf("ss0: %s\n", hex.EncodeToString(ss0[:]))
		t.Logf("ss0: %s\n", hex.EncodeToString(ss1[:]))

		t.FailNow()
	}
}

func TestHexSaveLoad(t *testing.T) {
	// create private key
	sk, err := ppk.NewPrivateKey()
	tests.AssertNil(t, err)

	// generate hex from key
	hexSk := sk.ToHex()

	// decode hex string into new key
	var newSk ppk.PrivateKey
	err = newSk.FromHex(hexSk)
	tests.AssertNil(t, err)

	// compare the two keys
	if !newSk.Equals(sk) {
		t.Fatalf("failed to match keys after hex save, load")
	}

	// apply the same for the public key
	pk := sk.PublicKey()
	hexPk := pk.ToHex()
	var newPk ppk.PublicKey
	err = newPk.FromHex(hexPk)
	tests.AssertNil(t, err)

	if !newPk.Equals(pk) {
		t.Fatalf("failed to match keys after hex save, load")
	}
}

func TestZeroingKey(t *testing.T) {
	// generate secret key
	sk, err := ppk.NewPrivateKey()
	tests.AssertNil(t, err)

	// zero key
	crypt.ZeroBytes(sk[:])

	// check if zero
	if !sk.IsZero() {
		t.Fatalf("failed to zero key")
	}

	// check key as hex string
	hex0 := "0000000000000000000000000000000000000000000000000000000000000000"
	if sk.ToHex() != hex0 {
		t.Fatalf("failed to match hex key to all 0 hex string")
	}

	// apply similar test for public key
	pk := sk.PublicKey()
	crypt.ZeroBytes(pk[:])
	if !pk.IsZero() {
		t.Fatalf("failed to zero key")
	}
}

func TestFromHexInvalidString(t *testing.T) {
	// test data
	invalidHexStrings := []string{
		"deadbeeff1f0",
		"contains invalid hex",
		"deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefcafe",
	}

	// test decoding the key and catch any missed errors
	var sk ppk.PrivateKey
	for _, str := range invalidHexStrings {
		err := sk.FromHex(str)
		if err == nil {
			t.Fatalf("decoded invalid key: %s\n", str)
		}
	}
}
