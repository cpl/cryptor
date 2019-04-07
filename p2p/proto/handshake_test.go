package proto_test

import (
	"bytes"
	"testing"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/proto"
)

func TestHandshake(t *testing.T) {
	// initializer
	isk, _ := ppk.NewPrivateKey()
	ihs := new(proto.Handshake)

	// receiver
	rsk, _ := ppk.NewPrivateKey()
	rpk := rsk.PublicKey()
	rhs := new(proto.Handshake)

	// protocol

	// I | initialize handshake, I knows the static pub key of R (rpk)
	imsg := ihs.Initialize(isk, rpk)

	// R | receive initializing message (imsg) and use own (rsk) to confirm
	//     at this point we extract the static pub key of I (ipkr)
	ipkr, err := rhs.Receive(imsg, rsk)
	if err != nil {
		t.Fatal(err)
	}

	// R | if ok, continue with response knowing (ipkr)
	rmsg := rhs.Respond(imsg.PlaintextUniquePublicKey, ipkr)

	// I | receive response, if ok, handshake is complete
	if err := ihs.Complete(rmsg, isk); err != nil {
		t.Fatal(err)
	}

	// at this point, both handshake hashes should match
	if !bytes.Equal(rhs.Hash[:], ihs.Hash[:]) {
		t.Fatal("failed to match handshake hashes")
	}

	// finalize handshakes and derive transport keys
	rsend, rrecv := rhs.Finalize()
	isend, irecv := ihs.Finalize()

	// check keys to not be zero
	var zeroKey [ppk.KeySize]byte
	if bytes.Equal(zeroKey[:], rsend[:]) || bytes.Equal(zeroKey[:], isend[:]) {

	}

	// compare keys
	if !bytes.Equal(rsend[:], irecv[:]) || !bytes.Equal(isend[:], rrecv[:]) {
		t.Fatal("failed to match transport keys")
	}
}
