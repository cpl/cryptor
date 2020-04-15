package crypt_test

import (
	"bytes"
	"math"
	"testing"

	"cpl.li/go/cryptor/internal/crypt"
	"cpl.li/go/cryptor/internal/crypt/hashing"

	"github.com/stretchr/testify/assert"
	chacha "golang.org/x/crypto/chacha20poly1305"
)

func TestRandomBytes(t *testing.T) {
	t.Parallel()

	for pow := 0.0; pow <= 14; pow++ {
		gen := int(math.Pow(2, pow))
		zero := make([]byte, gen)

		randBytes := crypt.RandomBytes(uint(gen))

		assert.Len(t, randBytes, gen, "failed to generate random bytes")
		assert.NotEqual(t, randBytes, zero)
	}
}

func TestRandomUint64(t *testing.T) {
	t.Parallel()
	size := 10

	values := make([]uint64, size)
	for _, v := range values {
		assert.Zero(t, v, "got unexpected 0 value")
	}

	for i := 0; i < size; i++ {
		values[i] = crypt.RandomUint64()
	}

	for _, v := range values {
		assert.NotZero(t, v, "got unexpected 0 value")
	}
}

func TestZeroBytes(t *testing.T) {
	t.Parallel()
	size := 10

	values := make([][]byte, size)
	for i := 0; i < size; i++ {
		values[i] = crypt.RandomBytes(1024)
	}

	crypt.ZeroBytes(values...)

	for _, v := range values {
		assert.Equal(t, make([]byte, len(v)), v, "found non-zero array")
	}
}

func TestAEADEncryption(t *testing.T) {
	t.Parallel()

	msg := []byte("We attack at dawn")
	key := crypt.RandomBytes(chacha.KeySize)

	cipher, err := chacha.New(key)
	assert.Nil(t, err)

	var hash hashing.HashSum

	nonce := crypt.RandomBytes(chacha.NonceSize)
	hashing.Hash(&hash, msg)

	ciphertext := cipher.Seal(nil, nonce, msg, hash[:])

	newCipher, err := chacha.New(key)
	assert.Nil(t, err)

	decrypted, err := newCipher.Open(nil, nonce, ciphertext, hash[:])
	assert.Nil(t, err)

	if !bytes.Equal(decrypted, msg) {
		t.Fatalf("message does not match decrypted message: %s : %s\n",
			string(decrypted), string(msg))
	}
}
