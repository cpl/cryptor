/*
Package mwords provides a list of mnemonic words which can be used for
translating binary/raw/... cryptographic information to a human readable format.
The word list used is as defined in BIP-39.
https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
*/
package mwords // import "cpl.li/go/cryptor/crypt/mwords"

import (
	"crypto/sha256"
	"errors"
	"math/big"
	"strings"
)

// bit mask for 11 least significant bits
var bits0to11 = big.NewInt(0x7FF)

// Count is the total number of words in the mnemonic list.
const Count = 2048

// MnemonicSentence is an array of mnemonic words, with extra utility methods
// on top.
type MnemonicSentence []string

// String will return all the words composing the mnemonic sentence as a single
// string of space seperated words.
func (ms MnemonicSentence) String() string {
	return strings.Join(ms, " ")
}

// IsValid checks each word in the sentence to be valid, if any is not, returns
// false. All words must be valid mnemonic words to return true. Another check
// is the number of words in the sentence.
func (ms MnemonicSentence) IsValid() bool {
	// validate word count
	if len(ms)%3 != 0 || len(ms) < 12 || len(ms) > 24 {
		return false
	}

	// validate words
	for _, word := range ms {
		if !IsValid(word) {
			return false
		}
	}
	return true
}

// IsValid checks if the given word is part of the mnemonic word list.
func IsValid(word string) bool {
	_, ok := mnemonicLookup[word]
	return ok
}

// ToMnemonic generates a BIP-39 mnemonic sentence that satisfies the given
// entropy length.
func ToMnemonic(entropy []byte) (MnemonicSentence, error) {
	// compute entropy bitcount
	bitsEntropy := len(entropy) * 8

	// validate entropy
	if !isValidEntropy(uint(bitsEntropy)) {
		return nil, errors.New("invalid entropy bit count")
	}

	// compute hash of entropy
	h := sha256.New()
	if _, err := h.Write(entropy); err != nil {
		return nil, err
	}
	hSum := h.Sum(nil)

	// compute checksum bitcount
	bitsChecksum := bitsEntropy / 32

	// convert entropy to big.Int, allows for easier bitwise operations
	bigEntropy := new(big.Int).SetBytes(entropy)

	// append first ENT/32 bits from the hash to entropy
	// where ENT = bitCount(entropy)
	bigEntropy.Lsh(bigEntropy, uint(bitsChecksum))
	bigEntropy.Or(bigEntropy, big.NewInt(int64(hSum[0]>>uint(8-bitsChecksum))))

	// allocate the number of mnemonic words soon to be generated
	wordCount := (bitsEntropy + bitsChecksum) / 11
	words := make(MnemonicSentence, wordCount)

	// finally we split all the bits into groups of 11 (because 2^11 = 2048)
	// exactly our word count, thus each set of 11 bits represents an index
	wordIndex := big.NewInt(0)
	for iter := wordCount - 1; iter >= 0; iter-- {
		// get least significant 11 bits
		wordIndex.And(bigEntropy, bits0to11)

		// shift out least significant 11 bits
		bigEntropy.Rsh(bigEntropy, 11)

		// convert the 11 bits to index
		words[iter] = mnemonicWords[wordIndex.Uint64()]
	}

	return words, nil
}

// TODO Write reverse function FromMnemonic(ms MnemonicSentence) ([]byte, error)
