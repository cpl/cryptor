/*
Package mwords provides a list of mnemonic words which can be used for
translating binary/raw/... cryptographic information to a human readable format.
The word list used is as defined in BIP-39.
https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
*/
package mwords // import "cpl.li/go/cryptor/crypt/mwords"

import (
	"strings"
)

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
// false. All words must be valid mnemonic words to return true.
func (ms MnemonicSentence) IsValid() bool {
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
