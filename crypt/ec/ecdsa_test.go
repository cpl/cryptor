package ec_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/thee-engineer/cryptor/crypt/ec"
)

func TestECDSAKeyGeneration(t *testing.T) {
	t.Parallel()

	prv, err := ec.GenerateKey()
	if err != nil {
		t.Error(err)
	}

	ecdsa := prv.Export()
	if ecdsa.D != prv.D {
		t.Error("ecdsa: exported key does not match ec key")
	}
}

func generateKeyParis() (*ec.PrivateKey, *ec.PrivateKey, error) {
	// Generate first key
	key0, err := ec.GenerateKey()
	if err != nil {
		return nil, nil, err
	}
	// Generate second key
	key1, err := ec.GenerateKey()
	if err != nil {
		return nil, nil, err
	}

	return key0, key1, nil
}

func generateDoubleSecret(key0, key1 *ec.PrivateKey) ([]byte, []byte, error) {
	// Generate shared secret from key 0
	sec0, err := key0.GenerateSecret(&key1.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	// Generate shared secret from key 1
	sec1, err := key1.GenerateSecret(&key0.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return sec0, sec1, nil
}

func TestECDSASecretGeneration(t *testing.T) {
	t.Parallel()

	// Generate two key pairs
	key0, key1, err := generateKeyParis()
	if err != nil {
		t.Error(err)
	}

	// Generate secrets
	sec0, sec1, err := generateDoubleSecret(key0, key1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(len(key0.D.Bytes()))
	fmt.Println(len(sec0))

	// Compare the two shared secrets
	if !bytes.Equal(sec0, sec1) {
		t.Error("ecdsa: shared secrets do not match")
	}
}

func TestECDSASecretGenerationError(t *testing.T) {
	t.Parallel()

	// Generate P521 go ecdsa key
	ecdsa521Key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Error(err)
	}

	// Import the key
	key0 := new(ec.PrivateKey)
	key0.Import(ecdsa521Key)

	// Create P256 ecdsa key
	key1, err := ec.GenerateKey()
	if err != nil {
		t.Error(err)
	}

	_, err = key0.GenerateSecret(&key1.PublicKey)
	if err == nil {
		t.Error("ecdsa: generated secret using invalid types")
	}

	_, err = key1.GenerateSecret(&key0.PublicKey)
	if err == nil {
		t.Error("ecdsa: generated secret using invalid types")
	}
}
