package ppk

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"log"
)

// KeySizeBits ...
const KeySizeBits = 4096

// NewKey ...
func NewKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, KeySizeBits)
	if err != nil {
		log.Panic(err)
	}
	if err := key.Validate(); err != nil {
		log.Panic(err)
	}
	key.Precompute()
	return key
}

// Encode ...
func Encode(data interface{}) []byte {
	buffer := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(data); err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

// DecodePrivate ...
func DecodePrivate(data []byte) *rsa.PrivateKey {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	key := new(rsa.PrivateKey)
	if err := decoder.Decode(key); err != nil {
		log.Panic(err)
	}
	return key
}

// DecodePublic ...
func DecodePublic(data []byte) *rsa.PublicKey {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	key := new(rsa.PublicKey)
	if err := decoder.Decode(key); err != nil {
		log.Panic(err)
	}
	return key
}
