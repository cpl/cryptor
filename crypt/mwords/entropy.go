package mwords

const (
	entropyMaxBits  = 256
	entropyMinBits  = 128
	entropyMultiple = 32
)

// check that the number of bits is within max, min and a multiple of 32
func isValidEntropy(bits uint) bool {
	if (bits%entropyMultiple) != 0 ||
		bits < entropyMinBits || bits > entropyMaxBits {

		return false
	}
	return true
}
