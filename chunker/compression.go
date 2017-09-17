package chunker

import (
	"bytes"
	"compress/gzip"
)

// Compress ...
func Compress(data []byte) (out bytes.Buffer, err error) {
	// Create gzip writer
	w := gzip.NewWriter(&out)
	defer w.Close()

	// Write data to output buffer
	_, err = w.Write(data)
	if err != nil {
		panic(err)
	}

	return out, nil
}

// Decompress ...
func Decompress(data bytes.Buffer) ([]byte, error) {
	var buffer bytes.Buffer

	// Create gzip reader
	r, err := gzip.NewReader(&data)
	if err != nil {
		return nil, nil
	}
	defer r.Close()

	buffer.ReadFrom(r)

	return buffer.Bytes(), nil
}
