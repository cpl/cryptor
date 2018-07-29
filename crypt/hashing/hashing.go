package hashing

import (
	"crypto"
	"hash"
	"log"

	"github.com/thee-engineer/cryptor/utils"
	"golang.org/x/crypto/blake2b"
)

// CryptoHash ...
const CryptoHash = crypto.BLAKE2b_256

// HashSize ...
const HashSize = blake2b.Size256

// NullHash ...
var NullHash = [HashSize]byte{0}

// HashFunction ...
func HashFunction() hash.Hash {
	h, err := blake2b.New256(nil)
	if err != nil {
		log.Panic(err)
	}
	return h
}

// Hash ...
func Hash(dataSet ...[]byte) []byte {
	h := HashFunction()
	for _, data := range dataSet {
		w, err := h.Write(data)
		utils.CheckErr(err)

		if w != len(data) {
			log.Panicf("blake2: write len %d does not match data len", w)
		}
	}

	return h.Sum(nil)
}

// Sum ...
func Sum(data []byte) []byte {
	h := HashFunction()
	w, err := h.Write(data)
	if err != nil {
		log.Panic(err)
	}
	if w != len(data) {
		log.Panicf("blake2: write len %d does not match data len", w)
	}

	return h.Sum(data)
}
