package crypt

import (
	"encoding/hex"
)

// Encode returns the given data as hex encoded []byte
func Encode(data []byte) []byte {
	out := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(out, data)
	return out
}

// Decode takes a hex encoded data and returns the decoded data
func Decode(data []byte) ([]byte, error) {
	out := make([]byte, hex.DecodedLen(len(data)))
	read, err := hex.Decode(out, data)
	if err != nil {
		return nil, err
	}
	return out[:read], nil
}

// EncodeString returns the given data as hex encoded string
func EncodeString(data []byte) string {
	return string(Encode(data))
}

// DecodeString takes a hex encoded string and returns the decoded []byte
func DecodeString(data string) ([]byte, error) {
	return Decode([]byte(data))
}
