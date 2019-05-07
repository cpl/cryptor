package mwords

import (
	"testing"

	"cpl.li/go/cryptor/tests"
)

func TestEntropyBits(t *testing.T) {
	t.Parallel()

	// invalid entropy bits
	for bits := uint(0); bits < entropyMinBits; bits++ {
		if isValidEntropy(bits) {
			t.Errorf("validated invalid number of bits %d\n", bits)
		}
	}
	for bits := uint(entropyMaxBits + 1); bits < entropyMaxBits*2; bits++ {
		if isValidEntropy(bits) {
			t.Errorf("validated invalid number of bits %d\n", bits)
		}
	}

	// valid range
	for bits := uint(entropyMinBits); bits < entropyMaxBits+1; bits++ {
		if isValidEntropy(bits) && bits%entropyMultiple != 0 {
			t.Errorf("validated invalid number of bits %d\n", bits)
		}
	}
}

func assertRandomWords(t *testing.T, num uint) {
	tests.AssertEqual(t, uint(len(RandomWords(num))), num, "invalid word count")
}

func TestRandomWords(t *testing.T) {
	for i := uint(0); i < 100; i++ {
		assertRandomWords(t, i)
	}
}
