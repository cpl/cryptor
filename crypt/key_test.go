package crypt

import "testing"

func TestKey(t *testing.T) {
	key := NewKey()
	NewKeyFromBytes(key.Bytes())
	NewKeyFromBytes(RandomData(AESKeySize))
	NewKeyFromString(key.String())
	NewKeyFromString("")
	key.Encode()
}
