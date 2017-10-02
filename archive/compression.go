package archive

import (
	"bytes"
	"compress/gzip"
)

// Compress takes an array of bytes and returns a byte buffer containing the
// gzip compressed data.
func Compress(data []byte) (out bytes.Buffer, err error) {
	// Create gzip writer
	gzipWriter := gzip.NewWriter(&out)
	defer gzipWriter.Close()

	// Write data to output buffer
	_, err = gzipWriter.Write(data)
	if err != nil {
		panic(err)
	}

	return out, nil
}

// Decompress takes a byte buffer containing gzip compressed data and returns
// the initial data as a byte array.
func Decompress(data bytes.Buffer) ([]byte, error) {
	var buffer bytes.Buffer

	// Create gzip reader
	gzipReader, err := gzip.NewReader(&data)
	if err != nil {
		return nil, nil
	}
	defer gzipReader.Close()

	buffer.ReadFrom(gzipReader)

	return buffer.Bytes(), nil
}
