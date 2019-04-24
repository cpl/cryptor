package noise_test

import (
	"bytes"
	"testing"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/tests"
)

func TestMarshalMessagesInitializer(t *testing.T) {
	t.Parallel()

	// generate test message
	secret, _ := ppk.NewPrivateKey()
	public := secret.PublicKey()
	msgI := new(noise.MessageInitializer)
	msgI.PlaintextUniquePublic = public

	// marshal into binary
	dataI, err := msgI.MarshalBinary()
	tests.AssertNil(t, err)

	// check size
	if size := len(dataI); size != noise.SizeMessageInitializer {
		t.Fatalf("invalid binary form size, expected %d, got %d\n",
			noise.SizeMessageInitializer, size)
	}

	newMsgI := new(noise.MessageInitializer)

	// unmarshal with invalid data
	if err := newMsgI.UnmarshalBinary(crypt.RandomBytes(50)); err == nil {
		t.Fatal("unmarshal invalid data")
	}

	// unmarshal
	if err := newMsgI.UnmarshalBinary(dataI); err != nil {
		t.Fatal(err)
	}

	// compare initial message with new message
	if !bytes.Equal(
		newMsgI.PlaintextUniquePublic[:], msgI.PlaintextUniquePublic[:]) {
		t.Fatal("failed to match keys")
	}
}

func TestMarshalMessagesResponder(t *testing.T) {
	t.Parallel()

	// generate test message
	secret, _ := ppk.NewPrivateKey()
	public := secret.PublicKey()
	msgR := new(noise.MessageResponder)
	msgR.PlaintextUniquePublic = public

	// marshal into binary
	dataR, err := msgR.MarshalBinary()
	tests.AssertNil(t, err)

	// check size
	if size := len(dataR); size != noise.SizeMessageResponder {
		t.Fatalf("invalid binary form size, expected %d, got %d\n",
			noise.SizeMessageResponder, size)
	}
	newMsgR := new(noise.MessageResponder)

	// unmarshal with invalid data
	if err := newMsgR.UnmarshalBinary(crypt.RandomBytes(50)); err == nil {
		t.Fatal("unmarshal invalid data")
	}

	// unmarshal
	if err := newMsgR.UnmarshalBinary(dataR); err != nil {
		t.Fatal(err)
	}

	// compare initial message with new message
	if !bytes.Equal(
		newMsgR.PlaintextUniquePublic[:], msgR.PlaintextUniquePublic[:]) {
		t.Fatal("failed to match keys")
	}
}
