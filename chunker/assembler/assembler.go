package assembler

import (
	"errors"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

// Assembler ...
type Assembler struct {
	data    []byte
	manager cachedb.Manager
}

// New ...
func New(manager cachedb.Manager) *Assembler {
	return &Assembler{
		manager: manager,
	}
}

// Unpack ...
func (a *Assembler) Unpack(hash []byte, password string) error {
	var data []byte
	var err error

	// Derive key
	key, err := aes.NewKeyFromPassword(password)
	if err != nil {
		return err
	}

	defer crypt.ZeroBytes(data, key[:])

	// While has next hash
	for {
		// Get chunk by hash
		data, err = a.manager.Get(hash)
		if err != nil {
			// TODO: Handle this case with P2P network and requests.
			return err
		}

		// Decrypt encrypted chunk
		data, err = aes.Decrypt(key, data)
		if err != nil {
			return err
		}

		// Extract chunk struct
		chk, err := chunker.ExtractChunk(data)
		if err != nil {
			return err
		}

		// Check chunk integrity
		if !chk.IsValid() {
			return errors.New("failed to validate chunk integrity")
		}

		// Check if chunk is the last chunk, exit loop and reconstruct data
		if chk.IsLast() {
			break
		}

		// Prepend decrypted data
		a.data = append(chk.Body, a.data...)

		// Update next hash and key
		hash = chk.Head.NextHash
		key = chk.Head.NextKey
	}

	return nil
}
