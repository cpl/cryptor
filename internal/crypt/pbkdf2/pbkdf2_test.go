package pbkdf2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cpl.li/go/cryptor/internal/crypt/pbkdf2"
	"cpl.li/go/cryptor/internal/crypt/ppk"
)

func TestPBKDF2(t *testing.T) {
	t.Parallel()

	expected :=
		"28df0b93627d5b50ed4fef574e774a00ac634cbd3395d0a57e769581e806f82f"
	expectedPub :=
		"a5f686a01f0327c2a1bce2d2ae01c4174d1637fd31a5a065d0b235ea37cc3d74"

	var (
		key ppk.PrivateKey
		pub ppk.PublicKey
	)

	key = pbkdf2.Key([]byte(password), []byte(salt))
	key.PublicKey(&pub)

	assert.Equal(t, len(key), ppk.KeySize, "invalid derived key length")

	assert.Equal(t, key.ToHex(), expected, "derived key does not match")

	assert.Equal(t, pub.ToHex(), expectedPub,
		"derived key public does not match")

	dKey := pbkdf2.Key([]byte(password), nil)
	if !key.Equals(dKey) {
		t.Fatal("default salt failed, keys don't match")
	}
}
