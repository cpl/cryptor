package pbkdf2_test

import (
	"testing"

	"cpl.li/go/cryptor/pkg/crypt/pbkdf2"
	"cpl.li/go/cryptor/pkg/crypt/ppk"
	"github.com/stretchr/testify/assert"
)

const password = "testing"
const salt = ".-_cryptor,$"

func TestPBKDF2(t *testing.T) {
	t.Parallel()

	expected :=
		"28df0b93627d5b50ed4fef574e774a00ac634cbd3395d0a57e769581e806f82f"
	expectedPub :=
		"a5f686a01f0327c2a1bce2d2ae01c4174d1637fd31a5a065d0b235ea37cc3d74"

	// derive key
	var key ppk.PrivateKey
	key = pbkdf2.Key([]byte(password), []byte(salt))

	// check len
	assert.Equal(t, len(key), ppk.KeySize, "invalid derived key length")

	// check expected key
	assert.Equal(t, key.ToHex(), expected, "derived key does not match")

	// check expected public key
	assert.Equal(t, key.PublicKey().ToHex(), expectedPub,
		"derived key public does not match")

	// check default salt is working
	dKey := pbkdf2.Key([]byte(password), nil)
	if !key.Equals(dKey) {
		t.Fatal("default salt failed, keys don't match")
	}
}

func BenchmarkPBKDF2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pbkdf2.Key([]byte(password), []byte(salt))
	}
}
