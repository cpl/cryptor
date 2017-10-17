package crypt

import (
	"encoding/hex"
)

// Encode returns the given msg as hex encoded []byte
func Encode(msg []byte) []byte {
	out := make([]byte, hex.EncodedLen(len(msg)))
	hex.Encode(out, msg)
	return out
}

// Decode takes a hex encoded msg and returns the decoded msg
func Decode(msg []byte) ([]byte, error) {
	out := make([]byte, hex.DecodedLen(len(msg)))
	read, err := hex.Decode(out, msg)
	if err != nil {
		return nil, err
	}
	return out[:read], nil
}

// EncodeString returns the given msg as hex encoded string
func EncodeString(msg []byte) string {
	return string(Encode(msg))
}

// DecodeString takes a hex encoded string and returns the decoded []byte
func DecodeString(msg string) ([]byte, error) {
	return Decode([]byte(msg))
}
