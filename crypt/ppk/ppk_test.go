package ppk_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"cpl.li/go/cryptor/crypt/mwords"

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

func TestToMnemonic(t *testing.T) {
	// generate key
	sk, err := ppk.NewPrivateKey()
	tests.AssertNil(t, err)

	// generate mnemonic and check size
	mnemonic := sk.ToMnemonic()
	if l := len(mnemonic); l != ppk.MnemonicSize {
		t.Fatalf("expected mnemonic len %d, got %d\n", ppk.MnemonicSize, l)
	}

	mnemonic = sk.PublicKey().ToMnemonic()
	if l := len(mnemonic); l != ppk.MnemonicSize {
		t.Fatalf("expected mnemonic len %d, got %d\n", ppk.MnemonicSize, l)
	}

	// generate mnemonic from zero key
	var zeroKey ppk.PrivateKey
	mnemonic = zeroKey.ToMnemonic()

	// check last word
	if mnemonic[ppk.MnemonicSize-1] != "art" {
		t.Fatalf("expected word \"art\", got \"%s\" instead\n",
			mnemonic[ppk.MnemonicSize-1])
	}

	// iterate all words but last
	for _, word := range mnemonic[:ppk.MnemonicSize-2] {
		if word != "abandon" {
			t.Errorf("expected word \"abandon\", got \"%s\" instead\n", word)
		}
	}
}

func TestFromMnemonic(t *testing.T) {
	// generate key
	sk, err := ppk.NewPrivateKey()
	tests.AssertNil(t, err)

	// generate mnemonics and check sizes
	mnemonicPrivate := sk.ToMnemonic()
	if l := len(mnemonicPrivate); l != ppk.MnemonicSize {
		t.Fatalf("expected mnemonic len %d, got %d\n", ppk.MnemonicSize, l)
	}
	mnemonicPublic := sk.PublicKey().ToMnemonic()
	if l := len(mnemonicPublic); l != ppk.MnemonicSize {
		t.Fatalf("expected mnemonic len %d, got %d\n", ppk.MnemonicSize, l)
	}

	// key generated from mnemonics
	var skDecoded ppk.PrivateKey
	var pkDecoded ppk.PublicKey

	// decode private key
	if err := skDecoded.FromMnemonic(mnemonicPrivate); err != nil {
		t.Fatal(err)
	}

	// decode public key
	if err := pkDecoded.FromMnemonic(mnemonicPublic); err != nil {
		t.Fatal(err)
	}

	// verify key integrity
	if !sk.Equals(skDecoded) {
		t.Errorf("private keys do not match")
	}
	if !sk.PublicKey().Equals(pkDecoded) {
		t.Errorf("public keys do not match")
	}
	if !skDecoded.PublicKey().Equals(pkDecoded) {
		t.Errorf("decoded public key error")
	}
}

func TestFromMnemonicZeroKey(t *testing.T) {
	// zero key
	var zeroKey ppk.PrivateKey

	// import mnemonic for zero key
	mnemonic, _ := mwords.MnemonicFromString(
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art")

	// generate random private key
	sk, _ := ppk.NewPrivateKey()

	// replace private key with key from mnemonic
	if err := sk.FromMnemonic(mnemonic); err != nil {
		t.Fatal(err)
	}

	// compare keys
	if !zeroKey.Equals(sk) {
		t.Fatal("zero keys do not match")
	}
}

func TestFromMnemonicInvalid(t *testing.T) {
	var key ppk.PrivateKey

	// empty mnemonic
	tests.AssertNotNil(t, key.FromMnemonic([]string{}), "empty mnemonic")

	// invalid mnemonic size
	tests.AssertNotNil(t,
		key.FromMnemonic(
			strings.Fields("legal year wave sausage worth useful legal winner thank yellow")), "invalid mnemonic size")

	// short mnemonic (half the wanted size)
	tests.AssertNotNil(t,
		key.FromMnemonic(
			strings.Fields("legal winner thank year wave sausage worth useful legal winner thank yellow")), "short mnemonic")

	// mnemonic contains invalid words
	tests.AssertNotNil(t,
		key.FromMnemonic(
			strings.Fields("beyond stage linux clip because twist token leaf atom foobarword genius food business side grid unable middle armed observe pair crouch tonight away coconut")), "invalid mnemonic")

	// mnemonic checksum is not valid
	tests.AssertNotNil(t,
		key.FromMnemonic(
			strings.Fields("zoo stage dog clip because twist token leaf atom about genius food business side grid unable middle armed observe pair crouch tonight away coconut")), "invalid mnemonic checksum")
}
