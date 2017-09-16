package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Encrypt ...
func Encrypt(key, msg []byte) ([]byte, error) {
	// Generate Cipher block
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Generate nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	// Generate GCM
	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	out := gcm.Seal(nil, nonce, msg, nil)

	return out, nil
}

// Decrypt ...
func Decrypt(key, msg, nonce []byte) ([]byte, error) {
	// Generate Cipher block
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := gcm.Open(nil, nonce, msg, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
