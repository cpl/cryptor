package noise

import (
	chacha "golang.org/x/crypto/chacha20poly1305"
)

var zeroNonce [chacha.NonceSize]byte

const (
	encryptedKeySize     = 48
	encryptedNothingSize = 16
)

// EncryptedKey ...
type EncryptedKey [encryptedKeySize]byte

// EncryptedNothing ...
type EncryptedNothing [encryptedNothingSize]byte
