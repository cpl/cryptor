package b64

import "encoding/base64"

// Encode ...
func Encode(data []byte) []byte {
	out := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(out, data)
	return out
}

// EncodeString ...
func EncodeString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Decode ...
func Decode(data []byte) ([]byte, error) {
	out := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(out, data)
	return out[:n], err
}

// DecodeString ...
func DecodeString(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
