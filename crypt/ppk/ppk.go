package ppk

import (
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"cpl.li/go/cryptor/crypt/mwords"
)

// KeySize is the size (in bytes) of both the public and private keys.
const KeySize = 32

// MnemonicSize is the number of words a key mnemonic will have. A key is 32
// bytes, so 256 bits, which BIP39 uses 24 words to represent.
const MnemonicSize = 24

type (
	// PrivateKey is a curve25519 point.
	PrivateKey [KeySize]byte

	// PublicKey is a curve25519 point, derived from the private key.
	PublicKey [KeySize]byte
)

func loadHexKey(src string, key []byte) error {
	// decode the string
	slice, err := hex.DecodeString(src)
	if err != nil {
		return err
	}

	// match with expected key size, return error if not matching
	if len(slice) != len(key) {
		return errors.New("hex string does not match key size")
	}

	// copy decoded string to key
	copy(key, slice)

	return nil
}

func loadMnemonicKey(mnemonic mwords.MnemonicSentence, key []byte) error {
	// validate mnemonic
	if len(mnemonic) != MnemonicSize {
		return errors.New("invalid mnemonic word count")
	}
	if !mnemonic.IsValid() {
		return errors.New("invalid mnemonic")
	}

	// decode and validate mnemonic
	b, err := mwords.EntropyFromMnemonic(mnemonic)
	if err != nil {
		return err
	}
	if len(b) != KeySize {
		return errors.New("invalid key size generated from mnemonic")
	}

	copy(key, b)

	return nil
}

// ToHex returns the secret key as a hex string.
func (sk PrivateKey) ToHex() string {
	return hex.EncodeToString(sk[:])
}

// ToHex returns the public key as a hex string.
func (pk PublicKey) ToHex() string {
	return hex.EncodeToString(pk[:])
}

// FromHex loads a hex string as a secret key.
func (sk *PrivateKey) FromHex(src string) error {
	return loadHexKey(src, sk[:])
}

// FromHex loads a hex string as a public key.
func (pk *PublicKey) FromHex(src string) error {
	return loadHexKey(src, pk[:])
}

// Equals compares two private keys, returns true if equal, false if not.
func (sk PrivateKey) Equals(to PrivateKey) bool {
	return subtle.ConstantTimeCompare(sk[:], to[:]) == 1
}

// Equals compares two public keys, returns true if equal, false if not.
func (pk PublicKey) Equals(to PublicKey) bool {
	return subtle.ConstantTimeCompare(pk[:], to[:]) == 1
}

// IsZero checks the key against a zero key, returns true if all key bytes are 0.
func (sk PrivateKey) IsZero() bool {
	var zero PrivateKey
	return sk.Equals(zero)
}

// IsZero checks the key against a zero key, returns true if all key bytes are 0.
func (pk PublicKey) IsZero() bool {
	var zero PublicKey
	return pk.Equals(zero)
}

// ToMnemonic exports the key as a 24 word mnemonic.
func (sk PrivateKey) ToMnemonic() mwords.MnemonicSentence {
	mnemonic, _ := mwords.EntropyToMnemonic(sk[:])
	return mnemonic
}

// ToMnemonic exports the key as a 24 word mnemonic.
func (pk PublicKey) ToMnemonic() mwords.MnemonicSentence {
	mnemonic, _ := mwords.EntropyToMnemonic(pk[:])
	return mnemonic
}

// FromMnemonic imports the key from a mnemonic sentence.
func (sk *PrivateKey) FromMnemonic(mnemonic mwords.MnemonicSentence) error {
	return loadMnemonicKey(mnemonic, sk[:])
}

// FromMnemonic imports the key from a mnemonic sentence.
func (pk *PublicKey) FromMnemonic(mnemonic mwords.MnemonicSentence) error {
	return loadMnemonicKey(mnemonic, pk[:])
}
