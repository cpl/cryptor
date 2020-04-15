package ppk

import (
	"crypto/subtle"
	"encoding/hex"

	"cpl.li/go/cryptor/internal/crypt/mwords"
)

type PublicKey [KeySize]byte

func (pk *PublicKey) ToHex() string {
	return hex.EncodeToString(pk[:])
}

func (pk *PublicKey) FromHex(src string) error {
	return loadHexKey(src, pk[:])
}

func (pk *PublicKey) Equals(to PublicKey) bool {
	return subtle.ConstantTimeCompare(pk[:], to[:]) == 1
}

func (pk *PublicKey) IsZero() bool {
	return pk.Equals(PublicKey{})
}

func (pk *PublicKey) ToMnemonic() mwords.MnemonicSentence {
	mnemonic, _ := mwords.EntropyToMnemonic(pk[:])
	return mnemonic
}

func (pk *PublicKey) FromMnemonic(mnemonic mwords.MnemonicSentence) error {
	return loadMnemonicKey(mnemonic, pk[:])
}
