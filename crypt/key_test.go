package crypt

import "testing"

func TestKey(t *testing.T) {
	t.Parallel()

	key := NewKey()
	NewKeyFromBytes(key.Bytes())
	NewKeyFromBytes(RandomData(AESKeySize))
	NewKeyFromString(key.String())
	NewKeyFromString("")
	key.Encode()
}
