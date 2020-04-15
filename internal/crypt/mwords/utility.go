package mwords

import (
	"math/big"
	"math/rand"
	"time"
)

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

var wordCountToEntropyBits = map[int]int{
	12: 128,
	15: 160,
	18: 192,
	21: 224,
	24: 256,
}

var checksumMaskMap = map[int]*big.Int{
	128: big.NewInt(0x0F),
	160: big.NewInt(0x1F),
	192: big.NewInt(0x3F),
	224: big.NewInt(0x7F),
	256: big.NewInt(0xFF),
}

func isValidEntropy(bits uint) bool {
	if (bits%entropyMultiple) != 0 ||
		bits < entropyMinBits || bits > entropyMaxBits {

		return false
	}
	return true
}

func RandomWords(n uint) []string {
	rand.Seed(time.Now().UTC().UnixNano())

	ret := make([]string, n)

	for n > 0 {
		ret[n-1] = mnemonicWords[rand.Int()%Count]
		n--
	}

	return ret
}
