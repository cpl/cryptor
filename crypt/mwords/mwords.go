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
	if len(ms)%sentenceMultiple != 0 ||
		len(ms) < sentenceMinWords || len(ms) > sentenceMaxWords {
		return false
	}

	// validate individual words
	for _, word := range ms {
		if !IsValidWord(word) {
			return false
		}
	}
	return true
}

// MnemonicFromString will split the given sentence, validate it and
// return a mnemonic sentence if possible, otherwise return an error.
func MnemonicFromString(sentence string) (MnemonicSentence, error) {
	// split given string
	var split MnemonicSentence = strings.Fields(sentence)

	// validate sentence and words
	if !split.IsValid() {
		return nil, errors.New("invalid mnemonic sentence")
	}

	return split, nil
}

// EntropyToMnemonic generates a BIP-39 mnemonic sentence that satisfies the given
// entropy length.
func EntropyToMnemonic(entropy []byte) (MnemonicSentence, error) {
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
	bitsChecksum := bitsEntropy / entropyMultiple

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

// EntropyFromString automatically performs the string to mnemonic
// sentence conversion and then applies the FromMnemonic function to get
// the entropy.
func EntropyFromString(sentence string) ([]byte, error) {
	// convert and validate string to mnemonic sentence
	ms, err := MnemonicFromString(sentence)
	if err != nil {
		return nil, err
	}

	// return entropy or error
	return EntropyFromMnemonic(ms)
}

// EntropyFromMnemonic takes a BIP-39 mnemonic sentence and returns the initial
// entropy used. If the sentence is invalid, an error is returned.
func EntropyFromMnemonic(ms MnemonicSentence) ([]byte, error) {
	// check if valid mnemonic
	if !ms.IsValid() {
		return nil, errors.New("invalid mnemonic sentence")
	}

	// compute bit counts
	bitsEntropy := wordCountToEntropyBits[len(ms)]
	bitsChecksum := bitsEntropy / entropyMultiple

	// use a big.Int to decode words, easier bitwise operations
	decoder := big.NewInt(0)

	// iterate words
	for _, word := range ms {
		// get word index
		index := mnemonicLookup[word]

		decoder.Lsh(decoder, 11)
		decoder.Or(decoder, big.NewInt(int64(index)))
	}

	// shift out checksum bits
	checksumBits := big.NewInt(0)
	checksumBits.And(decoder, checksumMaskMap[bitsEntropy])
	decoder.Rsh(decoder, uint(bitsChecksum))

	// get byte array
	decoded := decoder.Bytes()

	// check if the byte array needs zero padding
	if len(decoded) != bitsEntropy/8 {
		padding := make([]byte, (bitsEntropy/8)-len(decoded))
		decoded = append(padding, decoded...)
	}

	// compute checksum
	h := sha256.New()
	if _, err := h.Write(decoded); err != nil {
		return nil, err
	}
	hSum := h.Sum(nil)

	// validate checksum bits
	if checksumBits.Cmp(big.NewInt(int64(hSum[0]>>uint(8-bitsChecksum)))) != 0 {
		return nil, errors.New("failed to validate checksum bits")
	}

	return decoded, nil
}
