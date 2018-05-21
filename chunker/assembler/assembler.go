package assembler

import (
	"github.com/thee-engineer/cryptor/cachedb"
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

// func unpack(e []byte, key aes.Key) (*chunker.Chunk, error) {
// 	d, err := aes.Decrypt(key, e)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }

// Unpack ...
func (a *Assembler) Unpack(tail []byte, password string) error {
	_, err := a.manager.Get(tail)
	if err != nil {
		return err
	}

	return nil
}
