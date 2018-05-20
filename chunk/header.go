package chunk

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

// HeaderSize ...
var HeaderSize = 4 + crypt.KeySize + hashing.HashFunction().Size()*2

type header struct {
	Hash     []byte  // Hash of the chunk content
	NextHash []byte  // Hash of the next chunk
	NextKey  aes.Key // Key for the next chunk
	Padding  uint32  // Byte size of the padding
}

func newHeader() *header {
	return &header{
		NextKey:  aes.NullKey,
		NextHash: make([]byte, hashing.HashFunction().Size()),
		Hash:     make([]byte, hashing.HashFunction().Size()),
		Padding:  0,
	}
}

func extractHeader(data []byte) (*header, error) {
	if len(data) < HeaderSize {
		return nil, errors.New("invalid header size")
	}

	var count int
	head := newHeader()
	copy(head.Hash, data[count:count+hashing.HashFunction().Size()])
	count += hashing.HashFunction().Size()
	copy(head.NextHash, data[count:count+hashing.HashFunction().Size()])
	count += hashing.HashFunction().Size()
	copy(head.NextKey[:], data[count:count+crypt.KeySize])
	count += crypt.KeySize
	head.Padding = binary.LittleEndian.Uint32(data[count:])

	return head, nil
}

// Bytes ...
func (h *header) Bytes() []byte {
	var count int
	data := make([]byte, HeaderSize)

	copy(data[count:count+hashing.HashFunction().Size()], h.Hash)
	count += hashing.HashFunction().Size()
	copy(data[count:count+hashing.HashFunction().Size()], h.NextHash)
	count += hashing.HashFunction().Size()
	copy(data[count:count+crypt.KeySize], h.NextKey.Bytes())
	count += crypt.KeySize

	// Convert uint32 to byte array
	uintConv := make([]byte, 4)
	binary.LittleEndian.PutUint32(uintConv, h.Padding)

	copy(data[count:], uintConv)

	return data
}

// Zero clears the chunk header
func (h *header) Zero() {
	crypt.ZeroBytes(h.Hash[:], h.NextKey[:], h.NextHash[:])
	h.Padding = 0
}

// Equal ...
func (h *header) Equal(other *header) bool {
	if !bytes.Equal(h.Hash, other.Hash) {
		return false
	}
	if !bytes.Equal(h.NextHash, other.NextHash) {
		return false
	}
	if !bytes.Equal(h.NextKey.Bytes(), other.NextKey.Bytes()) {
		return false
	}
	if h.Padding != other.Padding {
		return false
	}
	return true
}
