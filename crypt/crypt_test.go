package crypt_test

import (
	"bytes"
	"math"
	"testing"

	"cpl.li/go/cryptor/crypt"
)

func TestRandomBytes(t *testing.T) {
	// go from 1 byte to 16k
	for pow := 0.0; pow <= 14; pow++ {
		gen := int(math.Pow(2, pow))
		if got := len(crypt.RandomBytes(uint(gen))); got != gen {
			t.Errorf("failed to generate %d random bytes, got %d\n", gen, got)
		}
	}
}

func assertZero(t *testing.T, data []byte) {
	if !bytes.Equal(data, make([]byte, len(data))) {
		t.Fatalf("failed to erase bytes, %v\n", data)
	}
}

func TestZero(t *testing.T) {
	// test data
	data0 := []byte("password1234")
	data1 := []byte("01234")
	data2 := crypt.RandomBytes(1024)
	data3 := crypt.RandomBytes(4096)

	// zero data
	crypt.ZeroBytes(data0, data1, data2)
	crypt.ZeroBytes(data3)

	// check all are zero
	assertZero(t, data0)
	assertZero(t, data1)
	assertZero(t, data2)
	assertZero(t, data3)
}

func assertHash(t *testing.T, data, expected string) {
	var sum crypt.Blake2sHash
	crypt.Hash(&sum, []byte(data))
	if sum.ToHex() != expected {
		t.Fatalf("non matching hashes, got %s for %s\n",
			sum.ToHex(), expected)
	}
}

func TestHash(t *testing.T) {
	assertHash(t, "",
		"69217a3079908094e11121d042354a7c1f55b6482ca1a51e1b250dfd1ed0eef9")
	assertHash(t, "Hello, World!",
		"ec9db904d636ef61f1421b2ba47112a4fa6b8964fd4a0a514834455c21df7812")
}
