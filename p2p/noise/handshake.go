/*
Package noise implements a custom version of the Noise Protocol Framework
handshake protocol mixed with concepts from WireGuard VPN.
*/
package noise // import "cpl.li/go/cryptor/p2p/noise"

import (
	"errors"
	"sync"

	"golang.org/x/crypto/blake2s"
	chacha "golang.org/x/crypto/chacha20poly1305"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/hkdf"
	"cpl.li/go/cryptor/crypt/ppk"
)

// utility byte array containing NonceSize zeroes
var zeroNonce [chacha.NonceSize]byte

// StateHandshake is an abstraction for the underlaying type used for representing
// the handshake state.
type StateHandshake byte

const (
	// StateEmpty can be on both initializer and responder sides. It's the
	// default starting state for newly created handshakes.
	StateEmpty StateHandshake = 0

	// StateInitialized is for a handshake marked by the initializer after it's
	// created using `noise.Initialize`.
	StateInitialized StateHandshake = 1

	// StateResponded is for a handshake marked by the responder after it
	// received an initialization message. The handshake must be generated
	// by using `noise.Respond`.
	StateResponded StateHandshake = 2

	// StateReceived is after the handshakes makes a round-trip from the
	// initializer to the responder and back, where it's validated.
	StateReceived StateHandshake = 3

	// StateSuccessful is marked after `Finalize` is called on StateResponded
	// for the responder or StateReceived for the initializer.
	StateSuccessful StateHandshake = 4
)

// Handshake contains all the information necessary for establishing a custom
// implementation of the Noise Protocol handshake.
type Handshake struct {
	// state or stage of the handshake
	state StateHandshake

	// for locking the handshake while performing any operation
	sync.RWMutex

	hash    crypt.Blake2sHash  // chaining hash
	c, t, k [blake2s.Size]byte // chaining keys

	presharedKey [ppk.KeySize]byte // optional pre-shared key

	tempKeys struct {
		secret ppk.PrivateKey // handshake unique private key
		public ppk.PublicKey  // handshake unique public key
	}
}

// generate temporary handshake keys
func (hs *Handshake) keygen() {
	hs.Lock()
	hs.tempKeys.secret, _ = ppk.NewPrivateKey()
	hs.tempKeys.public = hs.tempKeys.secret.PublicKey()
	hs.Unlock()
}

// State returns the handshake state.
func (hs *Handshake) State() StateHandshake {
	hs.RLock()
	defer hs.RUnlock()
	return hs.state
}

// PublicKey returns the temporary handshake public key generated at creation.
func (hs *Handshake) PublicKey() ppk.PublicKey {
	return hs.tempKeys.public
}

// Initialize will prepare the handshake structure for exchanging keys with a
// foreign node (with a known public key, rSPub). This method will return the
// initializers static public key (iSPub) encrypted. The receiving node should
// be able to compute the same encryption key and validate the encISPub.
func Initialize(iSPub, rSPub ppk.PublicKey) (
	hs *Handshake, msg *MessageInitializer) {
	// create handshake
	hs = new(Handshake)

	// generate handshake keys
	hs.keygen()

	// hash responder public key
	crypt.Hash(&hs.hash, rSPub[:])

	// derive chaining key
	hkdf.HKDF(hs.hash[:], hs.tempKeys.public[:], &hs.c)

	// compute shared secret
	ss := hs.tempKeys.secret.SharedSecret(rSPub)

	// derive chaining key and encryption key
	hkdf.HKDF(hs.c[:], ss[:], &hs.c, &hs.k)

	// prep message
	msg = new(MessageInitializer)
	msg.PlaintextUniquePublic = hs.tempKeys.public

	// encrypt static initializer public key
	cipher, _ := chacha.New(hs.k[:])
	cipher.Seal(msg.EncryptedInitializerStaticPublicKey[:0],
		zeroNonce[:], iSPub[:], hs.hash[:])

	// mark as initialized
	hs.state = StateInitialized

	return
}

// Respond will expect an already initialized handshake encISPub and the session
// temporary public key for validation and recreating the initializers handshake
// state (chaining keys, hash, etc). The responder private key (rSec) and the
// session temporary public key generated by the initializer are also required.
// This method will return both a Handshake and the response for the initializer.
func Respond(msg *MessageInitializer, rSec ppk.PrivateKey) (
	hs *Handshake, iSPub ppk.PublicKey, rmsg *MessageResponder, err error) {
	if msg == nil {
		return nil, iSPub, nil, errors.New("got nil message")
	}

	// create new handshake
	hs = new(Handshake)

	// hash responder public key
	rSPub := rSec.PublicKey()
	crypt.Hash(&hs.hash, rSPub[:])

	// derive chaining key
	hkdf.HKDF(hs.hash[:], msg.PlaintextUniquePublic[:], &hs.c)

	// compute shared secret
	ss := rSec.SharedSecret(msg.PlaintextUniquePublic)

	// derive chaining key and encryption key
	hkdf.HKDF(hs.c[:], ss[:], &hs.c, &hs.k)

	// decrypt static initializer public key
	cipher, _ := chacha.New(hs.k[:])
	_, err = cipher.Open(iSPub[:0], zeroNonce[:],
		msg.EncryptedInitializerStaticPublicKey[:], hs.hash[:])
	if err != nil {
		return nil, iSPub, rmsg, err
	}

	// generate handshake keys
	hs.keygen()

	// prep message
	rmsg = new(MessageResponder)
	rmsg.PlaintextUniquePublic = hs.tempKeys.public

	// update chaining key
	hkdf.HKDF(hs.c[:], hs.tempKeys.public[:], &hs.c)

	// update hash
	crypt.Hash(&hs.hash, hs.hash[:], hs.tempKeys.public[:])

	// update chaining key from shared secrets
	ss = hs.tempKeys.secret.SharedSecret(msg.PlaintextUniquePublic)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c)
	ss = hs.tempKeys.secret.SharedSecret(iSPub)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c)

	// derive chaining key, hash update key and encryption key
	// using pre-shared key in the mix
	hkdf.HKDF(hs.c[:], hs.presharedKey[:], &hs.c, &hs.t, &hs.k)

	// update hash
	crypt.Hash(&hs.hash, hs.hash[:], hs.t[:])

	// encrypt nothing
	cipher, _ = chacha.New(hs.k[:])
	cipher.Seal(rmsg.EncryptedNothing[:0], zeroNonce[:], nil, hs.hash[:])

	// update hash
	crypt.Hash(&hs.hash, rmsg.EncryptedNothing[:])

	// mark as responded
	hs.state = StateResponded

	return
}

// Receive is the last action in the handshake protocol before being able to
// generate transport keys on the initializer side.
func (hs *Handshake) Receive(msg *MessageResponder, iSec ppk.PrivateKey) error {
	if msg == nil {
		return errors.New("got nil message")
	}

	// lock handshake
	hs.Lock()
	defer hs.Unlock()

	// check handshake state
	if hs.state != StateInitialized {
		return errors.New("invalid handshake state")
	}

	// save current handshake state
	hsSnapshot := new(Handshake)
	hsSnapshot.hash = hs.hash
	hsSnapshot.c = hs.c
	hsSnapshot.k = hs.k
	hsSnapshot.t = hs.t

	// update chaining key
	hkdf.HKDF(hs.c[:], msg.PlaintextUniquePublic[:], &hs.c)

	// update hash
	crypt.Hash(&hs.hash, hs.hash[:], msg.PlaintextUniquePublic[:])

	// update chaining key from shared secrets
	var ss [ppk.KeySize]byte
	ss = hs.tempKeys.secret.SharedSecret(msg.PlaintextUniquePublic)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c)
	ss = iSec.SharedSecret(msg.PlaintextUniquePublic)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c)

	// derive chaining key, hash update key and encryption key
	// using pre-shared key in the mix
	hkdf.HKDF(hs.c[:], hs.presharedKey[:], &hs.c, &hs.t, &hs.k)

	// update hash
	crypt.Hash(&hs.hash, hs.hash[:], hs.t[:])

	// decrypt nothing for validation
	cipher, _ := chacha.New(hs.k[:])
	_, err := cipher.Open(nil, zeroNonce[:],
		msg.EncryptedNothing[:], hs.hash[:])
	if err != nil {
		// return to initial state
		hs.hash = hsSnapshot.hash
		hs.c = hsSnapshot.c
		hs.k = hsSnapshot.k
		hs.t = hsSnapshot.t

		return err
	}

	// update hash
	crypt.Hash(&hs.hash, msg.EncryptedNothing[:])

	// mark as received
	hs.state = StateReceived

	return nil
}

// Finalize can be called after the handshake protocol is completed either the
// initializer or receiver. In this step the transport keys are derived from
// the final shared state of the handshake, while other unused fields are
// zeroed. The keypair generated on each side of the protocol must matchs.
// initializerRecv = responderSend and initializerSend = responderRecv
func (hs *Handshake) Finalize() (send, recv [ppk.KeySize]byte, err error) {
	// lock handshake
	hs.Lock()
	defer hs.Unlock()

	// check for initializer or receiver mode
	// default case send, recv keys are all 0
	switch hs.state {
	// initializer
	case StateReceived:
		hkdf.HKDF(hs.c[:], nil, &send, &recv)
		hs.state = StateSuccessful
	// receiver
	case StateResponded:
		hkdf.HKDF(hs.c[:], nil, &recv, &send)
		hs.state = StateSuccessful
	// invalid handshake state
	default:
		return send, recv, errors.New("invalid handshake state")
	}

	// zero unused data
	crypt.ZeroBytes(hs.k[:], hs.c[:], hs.t[:],
		hs.tempKeys.public[:], hs.tempKeys.secret[:])

	// return transport keypair
	return send, recv, nil
}
