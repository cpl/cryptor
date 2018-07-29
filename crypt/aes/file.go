package aes

import (
	"io/ioutil"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/scrypt"
)

// EncryptFiles is used to encrypt local config files and the caches. For
// encrypting chunks, the aes.Encrypt function in combination with compression
// must be used.
func EncryptFiles(password string, files ...string) error {
	// Generate key from password and random salt
	keyBytes, salt := scrypt.RandomSalt(password)
	defer crypt.ZeroBytes(keyBytes, salt)
	key, err := NewKeyFromBytes(keyBytes)
	if err != nil {
		return err
	}

	// Iterate each file
	for _, file := range files {
		// Read file contents
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		defer crypt.ZeroBytes(data)

		// Encrypt with key
		eData, err := Encrypt(key, data)
		if err != nil {
			return err
		}

		// Append salt to file
		eData = append(eData, salt...)

		// Write encrypted data
		if err := ioutil.WriteFile(file, eData, 0666); err != nil {
			return err
		}
	}

	return nil
}

// DecryptFiles is used for local config and cache files.
func DecryptFiles(password string, files ...string) error {
	// Iterate each file
	for _, file := range files {
		// Read encrypted file
		eData, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		// Extract salt and re-generate key from password
		keyBytes := scrypt.Scrypt(password, eData[len(eData)-scrypt.SaltSize:])
		defer crypt.ZeroBytes(keyBytes)
		key, err := NewKeyFromBytes(keyBytes)
		if err != nil {
			return err
		}

		// Decrypt ignoring salt
		data, err := Decrypt(key, eData[:len(eData)-scrypt.SaltSize])
		if err != nil {
			return err
		}

		// Write decrypted file
		if err := ioutil.WriteFile(file, data, 0666); err != nil {
			return err
		}
	}

	return nil
}
