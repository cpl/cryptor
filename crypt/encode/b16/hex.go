package b16

import "encoding/hex"

// Encode ...
func Encode(data []byte) []byte {
	out := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(out, data)
	return out
}

// EncodeString ...
func EncodeString(data []byte) string {
	return string(Encode(data))
}

// Decode ...
func Decode(data []byte) ([]byte, error) {
	out := make([]byte, hex.DecodedLen(len(data)))
	read, err := hex.Decode(out, data)
	if err != nil {
		return nil, err
	}
	return out[:read], nil
}

// DecodeString ...
func DecodeString(data string) ([]byte, error) {
	return Decode([]byte(data))
}
