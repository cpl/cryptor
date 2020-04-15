package mwords

import (
	"crypto/sha256"
	"errors"
	"math/big"
	"strings"
)

type MnemonicSentence []string

func (ms MnemonicSentence) String() string {
	return strings.Join(ms, " ")
}

func (ms MnemonicSentence) IsValid() bool {
	if len(ms)%sentenceMultiple != 0 ||
		len(ms) < sentenceMinWords || len(ms) > sentenceMaxWords {
		return false
	}

	for _, word := range ms {
		if !IsValidWord(word) {
			return false
		}
	}
	return true
}

func MnemonicFromString(sentence string) (MnemonicSentence, error) {
	var split MnemonicSentence = strings.Fields(sentence)

	if !split.IsValid() {
		return nil, errors.New("invalid mnemonic sentence")
	}

	return split, nil
}

func EntropyToMnemonic(entropy []byte) (MnemonicSentence, error) {
	bitsEntropy := len(entropy) * 8

	if !isValidEntropy(uint(bitsEntropy)) {
		return nil, errors.New("invalid entropy bit count")
	}

	h := sha256.New()
	if _, err := h.Write(entropy); err != nil {
		return nil, err
	}
	hSum := h.Sum(nil)

	bitsChecksum := bitsEntropy / entropyMultiple

	bigEntropy := new(big.Int).SetBytes(entropy)

	bigEntropy.Lsh(bigEntropy, uint(bitsChecksum))
	bigEntropy.Or(bigEntropy, big.NewInt(int64(hSum[0]>>uint(8-bitsChecksum))))

	wordCount := (bitsEntropy + bitsChecksum) / 11
	words := make(MnemonicSentence, wordCount)

	wordIndex := big.NewInt(0)
	for iter := wordCount - 1; iter >= 0; iter-- {
		wordIndex.And(bigEntropy, bits0to11)
		bigEntropy.Rsh(bigEntropy, 11)
		words[iter] = mnemonicWords[wordIndex.Uint64()]
	}

	return words, nil
}

func EntropyFromString(sentence string) ([]byte, error) {
	ms, err := MnemonicFromString(sentence)
	if err != nil {
		return nil, err
	}

	return EntropyFromMnemonic(ms)
}

func EntropyFromMnemonic(ms MnemonicSentence) ([]byte, error) {
	if !ms.IsValid() {
		return nil, errors.New("invalid mnemonic sentence")
	}

	bitsEntropy := wordCountToEntropyBits[len(ms)]
	bitsChecksum := bitsEntropy / entropyMultiple

	decoder := big.NewInt(0)

	for _, word := range ms {
		index := mnemonicLookup[word]
		decoder.Lsh(decoder, 11)
		decoder.Or(decoder, big.NewInt(int64(index)))
	}

	checksumBits := big.NewInt(0)
	checksumBits.And(decoder, checksumMaskMap[bitsEntropy])
	decoder.Rsh(decoder, uint(bitsChecksum))

	decoded := decoder.Bytes()

	if len(decoded) != bitsEntropy/8 {
		padding := make([]byte, (bitsEntropy/8)-len(decoded))
		decoded = append(padding, decoded...)
	}

	h := sha256.New()
	if _, err := h.Write(decoded); err != nil {
		return nil, err
	}
	hSum := h.Sum(nil)

	if checksumBits.Cmp(big.NewInt(int64(hSum[0]>>uint(8-bitsChecksum)))) != 0 {
		return nil, errors.New("failed to validate checksum bits")
	}

	return decoded, nil
}
