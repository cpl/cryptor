package crypt_test

import (
	"bytes"
	"crypto/rand"
	"math"
	"testing"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/tests"

	chacha "golang.org/x/crypto/chacha20poly1305"
)

func TestRandomBytes(t *testing.T) {
	// go from 1 byte to 16k
	for pow := 0.0; pow <= 14; pow++ {
		gen := int(math.Pow(2, pow))
		if got := len(crypt.RandomBytes(uint(gen))); got != gen {
			t.Errorf("failed to generate %d random bytes, got %d\n", gen, got)
		}
	}
}

func TestRandomUint64(t *testing.T) {
	var val uint64
	val = crypt.RandomUint64()
	if val == 0 {
		t.Fatal("unexpected 0 value")
	}
}

func assertZero(t *testing.T, data []byte) {
	if !bytes.Equal(data, make([]byte, len(data))) {
		t.Fatalf("failed to erase bytes, %v\n", data)
	}
}

func TestZero(t *testing.T) {
	// test data
	data0 := []byte("password1234")
	data1 := []byte("01234")
	data2 := crypt.RandomBytes(1024)
	data3 := crypt.RandomBytes(4096)

	var data4 [40]byte
	rand.Read(data4[:])

	var data5 []byte
	data5 = make([]byte, 50)
	rand.Read(data5)

	// zero data
	crypt.ZeroBytes(data0, data1, data2)
	crypt.ZeroBytes(data3)
	crypt.ZeroBytes(data4[:])

	crypt.ZeroBytes(data5[:20])
	crypt.ZeroBytes(data5[20:])

	// check all are zero
	assertZero(t, data0)
	assertZero(t, data1)
	assertZero(t, data2)
	assertZero(t, data3)
	assertZero(t, data4[:])
	assertZero(t, data5)
}

func assertHash(t *testing.T, data, expected string) {
	var sum crypt.Blake2sHash
	crypt.Hash(&sum, []byte(data))
	if sum.ToHex() != expected {
		t.Fatalf("non matching hashes, got %s for %s\n",
			sum.ToHex(), expected)
	}
}

func TestHash(t *testing.T) {
	assertHash(t, "",
		"69217a3079908094e11121d042354a7c1f55b6482ca1a51e1b250dfd1ed0eef9")
	assertHash(t, "Hello, World!",
		"ec9db904d636ef61f1421b2ba47112a4fa6b8964fd4a0a514834455c21df7812")
}

func TestHashNil(t *testing.T) {
	data := []byte("Hello, World!")

	// request hash to be returned not passed in first argument
	hash0 := crypt.Hash(nil, data)
	if hash0.ToHex() !=
		"ec9db904d636ef61f1421b2ba47112a4fa6b8964fd4a0a514834455c21df7812" {
		t.Fatalf("failed to match hash from nil request\n")
	}

	// request hash to be passed in as argument
	var sum crypt.Blake2sHash
	hash1 := crypt.Hash(&sum, data)
	if hash1 != nil {
		t.Fatalf("hash function returned non-nil byte array %v\n", hash1[:])
	}

	if sum != *hash0 {
		t.Fatalf("failed to match hashes\n")
	}
}

func TestAEADEncryption(t *testing.T) {
	msg := []byte("We attack at dawn")
	key := crypt.RandomBytes(chacha.KeySize)

	// create cypher with random key
	cipher, err := chacha.New(key)
	tests.AssertNil(t, err)

	// generate nonce and hash
	nonce := crypt.RandomBytes(chacha.NonceSize)
	hash := crypt.Hash(nil, msg)[:]

	// encrypt message with nonce and hash for message auth
	ciphertext := cipher.Seal(nil, nonce, msg, hash)

	// generate new cipher with the same key
	newCipher, err := chacha.New(key)
	tests.AssertNil(t, err)

	// decrypt message
	decrypted, err := newCipher.Open(nil, nonce, ciphertext, hash)
	tests.AssertNil(t, err)

	// compare messages
	if !bytes.Equal(decrypted, msg) {
		t.Fatalf("message does not match decrypted message: %s : %s\n",
			string(decrypted), string(msg))
	}
}
