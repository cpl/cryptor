package chunker

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

func fakeHeader() *header {
	return &header{
		Hash:     crypt.RandomData(hashing.HashSize),
		NextHash: crypt.RandomData(hashing.HashSize),
		NextKey:  aes.NewKey(),
		Padding:  100,
	}
}

func TestHeader(t *testing.T) {
	t.Parallel()

	head := fakeHeader()

	extractedHead, err := extractHeader(head.Bytes())
	if err != nil {
		t.Error(err)
	}

	if !head.Equal(extractedHead) {
		t.Errorf("headers are not equal")
	}

	if !bytes.Equal(extractedHead.Bytes(), head.Bytes()) {
		t.Errorf("failed to match header bytes")
	}

	extractedHead.Zero()
	for _, b := range extractedHead.Bytes() {
		if b != 0 {
			t.Errorf("found non zero byte")
		}
	}
}

func TestHeaderExtract(t *testing.T) {
	t.Parallel()

	head := fakeHeader()
	headerBytes := head.Bytes()
	extractedHeader, err := extractHeader(headerBytes)
	if err != nil {
		t.Error(err)
	}
	if !head.Equal(extractedHeader) {
		t.Errorf("headers do not match")
	}
}

func TestHeaderEqual(t *testing.T) {
	t.Parallel()

	head0 := fakeHeader()
	head1 := fakeHeader()
	head1.Padding = 101

	if bytes.Equal(head0.Bytes(), head1.Bytes()) {
		t.FailNow()
	}
	if head0.Equal(head1) {
		t.FailNow()
	}

	head1.Hash = head0.Hash
	if head0.Equal(head1) {
		t.FailNow()
	}
	head1.NextHash = head0.NextHash
	if head0.Equal(head1) {
		t.FailNow()
	}
	head1.NextKey = head0.NextKey
	if head0.Equal(head1) {
		t.FailNow()
	}
	head1.Padding = head0.Padding
	if !head0.Equal(head1) {
		t.Errorf("failed to match equal headers")
	}
}

func TestHeaderExtractErr(t *testing.T) {
	t.Parallel()
	if _, err := extractHeader(crypt.RandomData(99)); err == nil {
		t.Errorf("extracted invalid header")
	}
}
