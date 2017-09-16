package chunker

import (
	"bytes"
	"compress/gzip"
)

func compress(data []byte) (buffer bytes.Buffer, err error) {
	// Create gzip writer
	w := gzip.NewWriter(&buffer)

	// Write data to output buffer
	_, err = w.Write(data)
	if err != nil {
		return buffer, err
	}

	// Close gzip writer
	err = w.Close()
	if err != nil {
		return buffer, err
	}

	return buffer, nil
}

func decompress(buffer *bytes.Buffer) ([]byte, error) {
	// Create gzip reader
	r, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}

	data := make([]byte, len(buffer.Bytes()))

	// Read data to output buffer
	_, err = r.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
