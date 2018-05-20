package chunker

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/common/con"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

func TestChunk(t *testing.T) {
	t.Parallel()

	c := newChunk(con.KB)
	data := crypt.RandomData(500)
	n, err := c.Write(data)
	if err != nil {
		t.Fatal(err)
	}
	if n != 500 {
		t.Errorf("did not write 500 bytes")
	}
	if !bytes.Equal(data, c.Body[0:500]) {
		t.Errorf("data does not match chunk body")
	}
	// if !bytes.Equal(c.Head.Hash, hashing.Hash(data)) {
	// 	t.Errorf("hashes do not match")
	// }
	n, err = c.Write(crypt.RandomData(1000))
	if err == nil || n != 0 {
		t.Errorf("wrote too much data")
	}
	// if !c.IsValid() {
	// 	t.Errorf("is valid returned false")
	// }
	n, err = c.Write(crypt.RandomData(500))
	if err != nil {
		t.Fatal(err)
	}
	// if c.Head.Padding != 0 {
	// 	t.Errorf("padding is not 0")
	// }
	key := aes.NewKey()
	nkey := aes.NewKey()
	_, err = c.Pack(key, nkey, crypt.RandomData(32))
	if err != nil {
		t.Error(err)
	}
}

func TestChunkEdge(t *testing.T) {
	t.Parallel()

	c := newChunk(con.KB)
	n, err := c.Write(crypt.RandomData(con.KB))
	if n != con.KB || err != nil {
		t.Errorf("failed to write to chunk")
	}

	_, err = c.Pack(aes.NewKey(), aes.NewKey(), crypt.RandomData(32))
	if err != nil {
		t.Error(err)
	}
}
