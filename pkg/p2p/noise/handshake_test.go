package noise_test

import (
	"bytes"
	"testing"

	"cpl.li/go/cryptor/crypt"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/tests"
)

func TestHandshake(t *testing.T) {
	t.Parallel()

	// initializer
	iSSecret, _ := ppk.NewPrivateKey()
	iSPublic := iSSecret.PublicKey()

	// responder
	rSSecret, _ := ppk.NewPrivateKey()
	rSPublic := rSSecret.PublicKey()

	// initializer begin protocol
	iHandshake, msgI := noise.Initialize(iSPublic, rSPublic)
	tests.AssertEqual(t, iHandshake.State(), noise.StateInitialized,
		"unexpected handshake state")

	// responder receives handshake data
	rHandshake, iSPublicR, msgR, err :=
		noise.Respond(msgI, rSSecret)
	tests.AssertNil(t, err)
	tests.AssertEqual(t, rHandshake.State(), noise.StateResponded,
		"unexpected handshake state")

	if !bytes.Equal(iSPublicR[:], iSPublic[:]) {
		t.Fatalf("failed to match initializer public key")
	}

	// responder sends response to initializer
	tests.AssertNil(t, iHandshake.Receive(msgR, iSSecret))
	tests.AssertEqual(t, iHandshake.State(), noise.StateReceived,
		"unexpected handshake state")

	// both handshakes can compute transport keys
	iSend, iRecv, err := iHandshake.Finalize()
	tests.AssertNil(t, err)
	tests.AssertEqual(t, iHandshake.State(), noise.StateSuccessful,
		"unexpected handshake state")

	var zeroPubKey ppk.PublicKey

	// compare pub unique key with zero key
	iPubKey := iHandshake.PublicKey()
	if !bytes.Equal(zeroPubKey[:], iPubKey[:]) {
		t.Error(iPubKey.ToHex())
		t.Fatal("handshake key is not zero after handshake")
	}

	rSend, rRecv, err := rHandshake.Finalize()
	tests.AssertNil(t, err)
	tests.AssertEqual(t, rHandshake.State(), noise.StateSuccessful,
		"unexpected handshake state")

	// compare pub unique key with zero key
	rPubKey := rHandshake.PublicKey()
	if !bytes.Equal(zeroPubKey[:], rPubKey[:]) {
		t.Error(iPubKey.ToHex())
		t.Fatal("handshake key is not zero after handshake")
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

func TestHandshakePublicKey(t *testing.T) {
	t.Parallel()

	var zeroPrivateKey ppk.PrivateKey

	// generate handshake
	newPrivateKey, _ := ppk.NewPrivateKey()
	hs, msg := noise.Initialize(
		zeroPrivateKey.PublicKey(), newPrivateKey.PublicKey())

	// make sure unique key being sent out is the same as handshake pub key
	hsPub := hs.PublicKey()
	if !bytes.Equal(msg.PlaintextUniquePublic[:], hsPub[:]) {
		t.Error("public keys do not match")
		t.Error(hsPub.ToHex())
		t.Error(msg.PlaintextUniquePublic.ToHex())
		t.FailNow()
	}
}

const (
	errStringInvalidHandshakeState = "invalid handshake state"
	errStringFailedAuth            = "chacha20poly1305: message authentication failed"
	errStringNilMessage            = "got nil message"
)

func TestInvalidHandshake(t *testing.T) {
	t.Parallel()

	var zeroPrivateKey ppk.PrivateKey

	// empty handshake
	hs := new(noise.Handshake)

	// attempt invalid operations on empty handshake
	if _, _, err := hs.Finalize(); err.Error() != errStringInvalidHandshakeState {
		t.Fatal("performed handshake Finalize on empty handshake, err:", err)
	}
	if err := hs.Receive(new(noise.MessageResponder), zeroPrivateKey); err.Error() != errStringInvalidHandshakeState {
		t.Fatal("performed handshake Receive on empty handshake, err:", err)
	}

	// generate real handshake
	newPrivateKey, _ := ppk.NewPrivateKey()
	hs, msg := noise.Initialize(
		zeroPrivateKey.PublicKey(), newPrivateKey.PublicKey())

	// attempt invalid operation
	if _, _, err := hs.Finalize(); err.Error() != "invalid handshake state" {
		t.Fatal("performed handshake Finalize on invalid handshake, err:", err)
	}

	// attempt with invalid private key
	if _, _, _, err := noise.Respond(msg, zeroPrivateKey); err.Error() != errStringFailedAuth {
		t.Fatal(err)
	}

	// attempt with nil message
	if _, _, _, err := noise.Respond(nil, newPrivateKey); err.Error() != errStringNilMessage {
		t.Fatal(err)
	}

	// copy original message
	msgCopy := *msg

	// attempt with invalid message pub key
	crypt.ZeroBytes(msgCopy.PlaintextUniquePublic[:])
	if _, _, _, err := noise.Respond(&msgCopy, newPrivateKey); err.Error() != errStringFailedAuth {
		t.Fatal(err)
	}

	// copy original message
	msgCopy = *msg

	// attempt with invalid message encrypted content
	crypt.ZeroBytes(msgCopy.EncryptedInitializerStaticPublicKey[:])
	if _, _, _, err := noise.Respond(&msgCopy, newPrivateKey); err.Error() != errStringFailedAuth {
		t.Fatal(err)
	}

	// attempt with valid message for the wrong peer
	wrongPrivateKey, _ := ppk.NewPrivateKey()
	_, wrongMsg := noise.Initialize(
		zeroPrivateKey.PublicKey(), wrongPrivateKey.PublicKey())
	if _, _, _, err := noise.Respond(wrongMsg, newPrivateKey); err.Error() != errStringFailedAuth {
		t.Fatal(err)
	}

	// respond to valid handshake
	rhs, _, rmsg, err := noise.Respond(msg, newPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	// attempt invalid operation
	if err := rhs.Receive(rmsg, newPrivateKey); err.Error() != errStringInvalidHandshakeState {
		t.Fatal("performed handshake Receive on invalid handshake, err:", err)
	}

	// attempt Receive with nil message
	if err := rhs.Receive(nil, newPrivateKey); err.Error() != errStringNilMessage {
		t.Fatal(err)
	}

	// finish handshake
	if _, _, err := rhs.Finalize(); err != nil {
		t.Fatal(err)
	}

	// attempt invalid operation
	if err := rhs.Receive(rmsg, newPrivateKey); err.Error() != errStringInvalidHandshakeState {
		t.Fatal("performed handshake Receive on invalid handshake, err:", err)
	}

	// attempt invalid operation
	if _, _, err := rhs.Finalize(); err.Error() != errStringInvalidHandshakeState {
		t.Fatal("performed handshake Finalize on invalid handshake, err:", err)
	}

	// invalid private key
	if err := hs.Receive(rmsg, newPrivateKey); err.Error() != errStringFailedAuth {
		t.Fatal("performed handshake Receive on invalid handshake, err:", err)
	}

	// nil message valid key
	if err := hs.Receive(nil, zeroPrivateKey); err.Error() != errStringNilMessage {
		t.Fatal("performed handshake Receive on invalid handshake, err:", err)
	}

	// copy original message
	rmsgCopy := *rmsg

	// attempt with invalid message pub key
	crypt.ZeroBytes(rmsgCopy.PlaintextUniquePublic[:])
	if err := hs.Receive(&rmsgCopy, newPrivateKey); err.Error() != errStringFailedAuth {
		t.Fatal(err)
	}

	// copy original message
	rmsgCopy = *rmsg

	// attempt with invalid message encrypted content
	crypt.ZeroBytes(rmsgCopy.EncryptedNothing[:])
	if err := hs.Receive(&rmsgCopy, newPrivateKey); err.Error() != errStringFailedAuth {
		t.Fatal(err)
	}

	// perform valid handshake
	tests.AssertNil(t, hs.Receive(rmsg, zeroPrivateKey))
}
