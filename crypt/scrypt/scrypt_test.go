package scrypt_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt/scrypt"
)

func TestScrypt(t *testing.T) {
	t.Parallel()

	key0 := scrypt.Scrypt("hello world", []byte("easy salt"))
	key1 := scrypt.Scrypt("hello world", []byte("easy salt"))

	if !bytes.Equal(key0, key1) {
		t.Error("scrypt: derived keys don't match")
	}
}

func TestScryptRandomSalt(t *testing.T) {
	t.Parallel()

	key0, salt := scrypt.RandomSalt("hello world")
	key1 := scrypt.Scrypt("hello world", salt)

	if !bytes.Equal(key0, key1) {
		t.Error("scrypt: derived keys don't match")
	}
}

func TestScryptAllRandom(t *testing.T) {
	t.Parallel()

	key0, pass, salt := scrypt.AllRandom()
	key1 := scrypt.Scrypt(string(pass), salt)

	if !bytes.Equal(key0, key1) {
		t.Error("scrypt: derived keys don't match")
	}
}

func TestScryptErrors(t *testing.T) {
	t.Parallel()

}
