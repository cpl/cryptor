package ec_test

import "testing"
import "github.com/thee-engineer/cryptor/crypt/ec"

func BenchmarkECDSAKeyGen(b *testing.B) {
	for count := 0; count < b.N; count++ {
		if _, err := ec.GenerateKey(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkECDSASecretGen(b *testing.B) {
	staticKey, err := ec.GenerateKey()
	if err != nil {
		b.Fatal(err)
	}
	pubStaticKey := staticKey.PublicKey

	keys := make([]*ec.PrivateKey, b.N)

	for count := 0; count < b.N; count++ {
		key, err := ec.GenerateKey()
		if err != nil {
			b.Fatal(err)
		}
		keys[count] = key
	}

	b.ResetTimer()

	for _, key := range keys {
		if _, err := key.GenerateSecret(&pubStaticKey); err != nil {
			b.Fatal(err)
		}
	}
}
