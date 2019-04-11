package noise_test

import (
	"bytes"
	"testing"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/noise"
)

func TestHandshake(t *testing.T) {
	// initializer
	iSSecret, _ := ppk.NewPrivateKey()
	iSPublic := iSSecret.PublicKey()

	// responder
	rSSecret, _ := ppk.NewPrivateKey()
	rSPublic := rSSecret.PublicKey()

	// initializer begin protocol
	iHandshake, encISPub := noise.Initialize(iSPublic, rSPublic)

	// responder receives handshake data
	rHandshake, iSPublicR, encNothing, err :=
		noise.Respond(iHandshake.PublicKey(), rSSecret, encISPub)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(iSPublicR[:], iSPublic[:]) {
		t.Fatalf("failed to match initializer public key")
	}

	// responder sends response to initializer
	if err := iHandshake.Receive(
		rHandshake.PublicKey(), iSSecret, encNothing); err != nil {
		t.Fatal(err)
	}

	// both handshakes can compute transport keys
	iSend, iRecv, err := iHandshake.Finalize()
	if err != nil {
		t.Fatal(err)
	}
	rSend, rRecv, err := rHandshake.Finalize()
	if err != nil {
		t.Fatal(err)
	}

	// check for zero key
	var zeroKey [ppk.KeySize]byte
	if bytes.Equal(zeroKey[:], rSend[:]) || bytes.Equal(zeroKey[:], iSend[:]) {
		t.Fatal("transport keys may be zero")
	}

	// compare keys
	if !bytes.Equal(rSend[:], iRecv[:]) || !bytes.Equal(iSend[:], rRecv[:]) {
		t.Fatal("failed to match transport keys")
	}
}

// TODO Write test cases for invalid states
// TODO Write test cases for MITM and other attacks

func TestHandshakeWithPackets(t *testing.T) {
	// initializer
	iSSecret, _ := ppk.NewPrivateKey()
	iSPublic := iSSecret.PublicKey()

	// responder
	rSSecret, _ := ppk.NewPrivateKey()
	rSPublic := rSSecret.PublicKey()

	// initializer begin protocol
	iHandshake, encISPub := noise.Initialize(iSPublic, rSPublic)

	// create packet

	// responder receives handshake data
	rHandshake, iSPublicR, encNothing, err :=
		noise.Respond(iHandshake.PublicKey(), rSSecret, encISPub)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(iSPublicR[:], iSPublic[:]) {
		t.Fatalf("failed to match initializer public key")
	}

	// responder sends response to initializer
	if err := iHandshake.Receive(
		rHandshake.PublicKey(), iSSecret, encNothing); err != nil {
		t.Fatal(err)
	}

	// both handshakes can compute transport keys
	iSend, iRecv, err := iHandshake.Finalize()
	if err != nil {
		t.Fatal(err)
	}
	rSend, rRecv, err := rHandshake.Finalize()
	if err != nil {
		t.Fatal(err)
	}

	// check for zero key
	var zeroKey [ppk.KeySize]byte
	if bytes.Equal(zeroKey[:], rSend[:]) || bytes.Equal(zeroKey[:], iSend[:]) {
		t.Fatal("transport keys may be zero")
	}

	// compare keys
	if !bytes.Equal(rSend[:], iRecv[:]) || !bytes.Equal(iSend[:], rRecv[:]) {
		t.Fatal("failed to match transport keys")
	}
}
