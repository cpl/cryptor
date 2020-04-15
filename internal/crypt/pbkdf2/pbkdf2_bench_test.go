package pbkdf2_test

import (
	"testing"

	"cpl.li/go/cryptor/internal/crypt/ppk"

	"cpl.li/go/cryptor/internal/crypt/pbkdf2"
)

const (
	password = "testing"
	salt     = ".-_cryptor,$"
)

func BenchmarkPBKDF2(b *testing.B) {
	var key ppk.PrivateKey
	for i := 0; i < b.N; i++ {
		key = pbkdf2.Key([]byte(password), []byte(salt))
	}
	key.IsZero()
}
