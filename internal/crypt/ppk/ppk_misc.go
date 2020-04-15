package ppk

import (
	"encoding/hex"
	"errors"

	"cpl.li/go/cryptor/internal/crypt/mwords"
)

const (
	// KeySize ...
	KeySize = 32

	// MnemonicSize ...
	MnemonicSize = 24
)

func loadHexKey(src string, key []byte) error {
	slice, err := hex.DecodeString(src)
	if err != nil {
		return err
	}

	if len(slice) != len(key) {
		return errors.New("hex string does not match key size")
	}

	copy(key, slice)

	return nil
}

func loadMnemonicKey(mnemonic mwords.MnemonicSentence, key []byte) error {
	if len(mnemonic) != MnemonicSize {
		return errors.New("invalid mnemonic word count")
	}
	if !mnemonic.IsValid() {
		return errors.New("invalid mnemonic")
	}

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
