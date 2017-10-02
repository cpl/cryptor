package crypt

import (
	"testing"
)

func TestKey(t *testing.T) {
	t.Parallel()

	key := NewKey()
	NewKeyFromBytes(key.Bytes())
	NewKeyFromBytes(RandomData(AESKeySize))
	NewKeyFromString(key.String())
	NewKeyFromString("")
	key.Encode()
}

func TestKeyFPassword(t *testing.T) {
	t.Parallel()

	// Test passwords
	var testValues = []string{
		"",
		"0123",
		"testinglargerpassword",
		"5rdZ<Q/{Uf6@Ed!~uF8T(9Ad8S<+{w9,",
		">~Ph]p`}$_QSW6#wrfo$/a<fH+PLz,Ycv;tY#Y4b&2uNwC7BKT~GZ<HsPXt(}"}

	// Expected test results for password key derivation
	var testResult = []string{
		"64a868d4b23af696d3734d0b814d04cdd1ac280128e97653a05f32b49c13a29a",
		"65e4e5bf4cb95801ead6d6dfb3367bfcd3c7cfebd62efaee312f2b6d40ff4f5c",
		"00c4ce73053141d6b53016dc02bc4efaa65f38ed457ed1f8b28eb65499b7f555",
		"bae68a56d1acea882202bcc3fe8857235bb410df38c9ae28e7f2e44100f9aced",
		"cd65b591e9d0c68359f344c9c693e75b51544aa4ba88b32d7555517c447427a2"}

	// Test all values
	for index, test := range testValues {
		if testResult[index] != string(NewKeyFromPassword(test).Encode()) {
			t.Error("key error: wrong key derivation")
		}
	}
}
