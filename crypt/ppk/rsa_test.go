package ppk_test

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/common/con"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/ppk"
	"github.com/thee-engineer/cryptor/utils"
)

func TestEncryption(t *testing.T) {
	t.Parallel()

	// Generate data and key
	data := crypt.RandomData(128)
	key := ppk.NewKey()

	// Encrypt data
	eData, err := ppk.Encrypt(&key.PublicKey, data)
	utils.CheckErrTest(err, t)

	// Sanity check
	if bytes.Equal(eData, data) {
		t.Errorf("ppk encrypt | data still matches")
	}

	// Decrypt data
	dData, err := ppk.Decrypt(key, eData)
	utils.CheckErrTest(err, t)

	// Check if data matches
	if !bytes.Equal(data, dData) {
		t.Errorf("ppk decrypt | data does not match")
	}
}

func TestEncryptionErrors(t *testing.T) {
	t.Parallel()

	// Generate data and key
	data := crypt.RandomData(con.KB)
	key := ppk.NewKey()

	_, err := ppk.Encrypt(&key.PublicKey, data)
	if err.Error() != "crypto/rsa: message too long for RSA public key size" {
		t.Error(err)
	}

	_, err = ppk.Decrypt(key, data)
	if err.Error() != "crypto/rsa: decryption error" {
		t.Error(err)
	}

	eData, err := ppk.Encrypt(&key.PublicKey, crypt.RandomData(10))
	utils.CheckErrTest(err, t)

	nKey := ppk.NewKey()
	_, err = ppk.Decrypt(nKey, eData)
	if err.Error() != "crypto/rsa: decryption error" {
		t.Error(err)
	}

	eData[0] += 1
	_, err = ppk.Decrypt(key, eData)
	if err.Error() != "crypto/rsa: decryption error" {
		t.Error(err)
	}
}

func TestSignature(t *testing.T) {
	t.Parallel()

	data := crypt.RandomData(con.KB)
	key := ppk.NewKey()

	signature, err := ppk.Sign(key, data)
	utils.CheckErrTest(err, t)

	if !ppk.Verify(&key.PublicKey, data, signature) {
		t.Errorf("ppk verify failed")
	}
}

func TestSignatureErrors(t *testing.T) {
	t.Parallel()

	data := crypt.RandomData(con.KB)
	key0 := ppk.NewKey()
	key1 := ppk.NewKey()

	signature, err := ppk.Sign(key0, data)
	utils.CheckErrTest(err, t)

	if ppk.Verify(&key1.PublicKey, data, signature) {
		t.Errorf("ppk verified wrong key")
	}
	if ppk.Verify(&key0.PublicKey, data, data) {
		t.Errorf("ppk verified with wrong format")
	}
	if ppk.Verify(&key0.PublicKey, data, crypt.RandomData(512)) {
		t.Errorf("ppk verified random signature")
	}
	signature[0] += 1
	if ppk.Verify(&key0.PublicKey, data, signature) {
		t.Errorf("ppk verified bad signature")
	}
}

func TestLargeEncrypt(t *testing.T) {
	t.Parallel()

	rsaKey := ppk.NewKey()
	aesKey := aes.NewKey()
	data := crypt.RandomData(con.MB)

	eData, err := aes.Encrypt(aesKey, data)
	utils.CheckErrTest(err, t)

	secretKey, err := ppk.Encrypt(&rsaKey.PublicKey, aesKey.Bytes())
	utils.CheckErrTest(err, t)

	decryptedKeyBytes, err := ppk.Decrypt(rsaKey, secretKey)
	utils.CheckErrTest(err, t)

	decryptedKey, err := aes.NewKeyFromBytes(decryptedKeyBytes)
	if err != nil {
		t.Error(err)
	}
	dData, err := aes.Decrypt(decryptedKey, eData)
	utils.CheckErrTest(err, t)

	if !bytes.Equal(data, dData) {
		t.Errorf("rsa aes failed, data mismatch")
	}
}
