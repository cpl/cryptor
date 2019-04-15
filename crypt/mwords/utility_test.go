package mwords

import (
	"testing"
)

func TestEntropyBits(t *testing.T) {
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
