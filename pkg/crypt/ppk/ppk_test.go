package ppk_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"cpl.li/go/cryptor/pkg/crypt"
	"cpl.li/go/cryptor/pkg/crypt/mwords"
	"cpl.li/go/cryptor/pkg/crypt/ppk"

	"github.com/stretchr/testify/assert"
)

func TestSharedSecretGeneration(t *testing.T) {
	t.Parallel()

	// generate two private keys
	sk0, err := ppk.NewPrivateKey()
	assert.Nil(t, err)
	sk1, err := ppk.NewPrivateKey()
	assert.Nil(t, err)

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
	t.Parallel()

	// create private key
	sk, err := ppk.NewPrivateKey()
	assert.Nil(t, err)

	// generate hex from key
	hexSk := sk.ToHex()

	// decode hex string into new key
	var newSk ppk.PrivateKey
	err = newSk.FromHex(hexSk)
	assert.Nil(t, err)

	// compare the two keys
	if !newSk.Equals(sk) {
		t.Fatalf("failed to match keys after hex save, load")
	}

	// apply the same for the public key
	pk := sk.PublicKey()
	hexPk := pk.ToHex()
	var newPk ppk.PublicKey
	err = newPk.FromHex(hexPk)
	assert.Nil(t, err)

	if !newPk.Equals(pk) {
		t.Fatalf("failed to match keys after hex save, load")
	}
}

func TestZeroingKey(t *testing.T) {
	t.Parallel()

	// generate secret key
	sk, err := ppk.NewPrivateKey()
	assert.Nil(t, err)

	// zero key
	crypt.ZeroBytes(sk[:])

	// check if zero
	if !sk.IsZero() {
		t.Fatalf("failed to zero key")
	}

	// check key as hex string
	assert.Equal(t, sk.ToHex(),
		"0000000000000000000000000000000000000000000000000000000000000000",
		"failed to match hex key to all 0 hex string")

	// apply similar test for public key
	pk := sk.PublicKey()
	crypt.ZeroBytes(pk[:])
	if !pk.IsZero() {
		t.Fatalf("failed to zero key")
	}
}

func TestFromHexInvalidString(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	// generate key
	sk, err := ppk.NewPrivateKey()
	assert.Nil(t, err)

	// generate mnemonic and check size
	mnemonic := sk.ToMnemonic()
	assert.Equal(t, len(mnemonic), ppk.MnemonicSize,
		"invalid mnemonic length")
	mnemonic = sk.PublicKey().ToMnemonic()
	assert.Equal(t, len(mnemonic), ppk.MnemonicSize,
		"invalid mnemonic length")

	// generate mnemonic from zero key
	var zeroKey ppk.PrivateKey
	mnemonic = zeroKey.ToMnemonic()

	// check last word
	assert.Equal(t, mnemonic[ppk.MnemonicSize-1], "art",
		"got unexpected word")

	// iterate all words but last
	for _, word := range mnemonic[:ppk.MnemonicSize-2] {
		assert.Equal(t, word, "abandon",
			"got unexpected word")
	}
}

func TestFromMnemonic(t *testing.T) {
	t.Parallel()

	// generate key
	sk, err := ppk.NewPrivateKey()
	assert.Nil(t, err)

	// generate mnemonics and check sizes
	mnemonicPrivate := sk.ToMnemonic()

	assert.Equal(t, len(mnemonicPrivate), ppk.MnemonicSize,
		"invalid mnemonic length")
	mnemonicPublic := sk.PublicKey().ToMnemonic()
	assert.Equal(t, len(mnemonicPublic), ppk.MnemonicSize,
		"invalid mnemonic length")

	// key generated from mnemonics
	var skDecoded ppk.PrivateKey
	var pkDecoded ppk.PublicKey

	// decode private key
	assert.Nil(t, skDecoded.FromMnemonic(mnemonicPrivate))

	// decode public key
	assert.Nil(t, pkDecoded.FromMnemonic(mnemonicPublic))

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
	t.Parallel()

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
	t.Parallel()

	var key ppk.PrivateKey

	// empty mnemonic
	assert.NotNil(t, key.FromMnemonic([]string{}), "empty mnemonic")

	// invalid mnemonic size
	assert.NotNil(t,
		key.FromMnemonic(
			strings.Fields("legal year wave sausage worth useful legal winner thank yellow")), "invalid mnemonic size")

	// short mnemonic (half the wanted size)
	assert.NotNil(t,
		key.FromMnemonic(
			strings.Fields("legal winner thank year wave sausage worth useful legal winner thank yellow")), "short mnemonic")

	// mnemonic contains invalid words
	assert.NotNil(t,
		key.FromMnemonic(
			strings.Fields("beyond stage linux clip because twist token leaf atom foobarword genius food business side grid unable middle armed observe pair crouch tonight away coconut")), "invalid mnemonic")

	// mnemonic checksum is not valid
	assert.NotNil(t,
		key.FromMnemonic(
			strings.Fields("zoo stage dog clip because twist token leaf atom about genius food business side grid unable middle armed observe pair crouch tonight away coconut")), "invalid mnemonic checksum")
}