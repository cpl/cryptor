package crypt

import (
	"bytes"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	data := RandomData(32)
	dData, err := Decode(Encode(data))
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(dData, data) {
		t.Error("Initial data does not match encoded->decoded data!")
	}
}

func TestStringEncodeDecode(t *testing.T) {
	t.Parallel()

	data := RandomData(32)
	dData, err := DecodeString(EncodeString(data))
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(dData, data) {
		t.Error("Initial data does not match encoded->decoded data!")
	}
}

func TestDecodeErrors(t *testing.T) {
	t.Parallel()

	if _, err := Decode([]byte{0, 1, 2, 3, 4}); err != nil {
		if err.Error() != "encoding/hex: odd length hex string" {
			t.Error(err)
		}
	}

	if _, err := Decode([]byte{0, 1, 2, 3}); err != nil {
		if err.Error() != "encoding/hex: invalid byte: U+0000" {
			t.Error(err)
		}
	}
}
