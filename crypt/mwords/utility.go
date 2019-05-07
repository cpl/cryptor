package mwords

import (
	"math/big"
	"math/rand"
	"time"
)

// bit mask for 11 least significant bits
var bits0to11 = big.NewInt(0x7FF)

const (
	entropyMaxBits  = 256
	entropyMinBits  = 128
	entropyMultiple = 32
)

const (
	sentenceMaxWords = 24
	sentenceMinWords = 12
	sentenceMultiple = 3
)

// quick translation between word count of a mnemonic sentence and the original
// entropy bit count, excluding the checksum bit count
var wordCountToEntropyBits = map[int]int{
	12: 128,
	15: 160,
	18: 192,
	21: 224,
	24: 256,
}

// a simple map between the entropy bit count and the mask required to extract
// the checksum bits
var checksumMaskMap = map[int]*big.Int{
	128: big.NewInt(0x0F),
	160: big.NewInt(0x1F),
	192: big.NewInt(0x3F),
	224: big.NewInt(0x7F),
	256: big.NewInt(0xFF),
}

// check that the number of bits is within max, min and a multiple of 32
func isValidEntropy(bits uint) bool {
	if (bits%entropyMultiple) != 0 ||
		bits < entropyMinBits || bits > entropyMaxBits {

		return false
	}
	return true
}

// RandomWords will return n random words from the wordlist.
func RandomWords(n uint) []string {
	// seed RNG
	rand.Seed(time.Now().UTC().UnixNano())

	// return slice
	ret := make([]string, n)

	// get n random words
	for n > 0 {
		ret[n-1] = mnemonicWords[rand.Int()%Count]
		n--
	}

	return ret
}
