package ppk

import (
	"crypto/subtle"
	"encoding/hex"

	"cpl.li/go/cryptor/internal/crypt/mwords"
)

type PrivateKey [KeySize]byte

func (sk *PrivateKey) ToHex() string {
	return hex.EncodeToString(sk[:])
}

func (sk *PrivateKey) FromHex(src string) error {
	return loadHexKey(src, sk[:])
}

func (sk *PrivateKey) Equals(to PrivateKey) bool {
	return subtle.ConstantTimeCompare(sk[:], to[:]) == 1
}

func (sk *PrivateKey) IsZero() bool {
	return sk.Equals(PrivateKey{})
}

func (sk *PrivateKey) ToMnemonic() mwords.MnemonicSentence {
	mnemonic, _ := mwords.EntropyToMnemonic(sk[:])
	return mnemonic
}

func (sk *PrivateKey) FromMnemonic(mnemonic mwords.MnemonicSentence) error {
	return loadMnemonicKey(mnemonic, sk[:])
}
