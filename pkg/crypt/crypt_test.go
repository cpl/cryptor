package crypt_test

import (
	"bytes"
	"crypto/rand"
	"math"
	"testing"

	"cpl.li/go/cryptor/pkg/crypt"
	"cpl.li/go/cryptor/pkg/crypt/hashing"

	"github.com/stretchr/testify/assert"
	chacha "golang.org/x/crypto/chacha20poly1305"
)

func TestRandomBytes(t *testing.T) {
	t.Parallel()

	// go from 1 byte to 16k
	for pow := 0.0; pow <= 14; pow++ {
		gen := int(math.Pow(2, pow))
		if got := len(crypt.RandomBytes(uint(gen))); got != gen {
			t.Errorf("failed to generate %d random bytes, got %d\n", gen, got)
		}
	}
}

func TestRandomUint64(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	// test data
	data0 := []byte("password1234")
	data1 := []byte("01234")
	data2 := crypt.RandomBytes(1024)
	data3 := crypt.RandomBytes(4096)

	var data4 [40]byte
	_, err := rand.Read(data4[:])
	assert.Nil(t, err)

	var data5 []byte
	data5 = make([]byte, 50)
	_, err = rand.Read(data5)
	assert.Nil(t, err)

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

func TestAEADEncryption(t *testing.T) {
	t.Parallel()

	msg := []byte("We attack at dawn")
	key := crypt.RandomBytes(chacha.KeySize)

	// create cypher with random key
	cipher, err := chacha.New(key)
	assert.Nil(t, err)

	// generate nonce and hash
	nonce := crypt.RandomBytes(chacha.NonceSize)
	hash := hashing.Hash(nil, msg)[:]

	// encrypt message with nonce and hash for message auth
	ciphertext := cipher.Seal(nil, nonce, msg, hash)

	// generate new cipher with the same key
	newCipher, err := chacha.New(key)
	assert.Nil(t, err)

	// decrypt message
	decrypted, err := newCipher.Open(nil, nonce, ciphertext, hash)
	assert.Nil(t, err)

	// compare messages
	if !bytes.Equal(decrypted, msg) {
		t.Fatalf("message does not match decrypted message: %s : %s\n",
			string(decrypted), string(msg))
	}
}
