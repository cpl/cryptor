package crypto

import (
	"crypto/sha256"
	"hash"
	"io"
	"os"
)

// SHA256File ...
func SHA256File(path string) (hash.Hash, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}

	return h, nil
}

// SHA256Data ...
func SHA256Data(data []byte) (hash.Hash, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}

	return h, nil
}
