package p2p

import (
	"sync"

	"cpl.li/go/cryptor/crypt/hkdf"
	"cpl.li/go/cryptor/crypt/ppk"
	"golang.org/x/crypto/blake2s"
	chacha "golang.org/x/crypto/chacha20poly1305"

	"cpl.li/go/cryptor/crypt"
)

// utility byte array containing NonceSize zeroes
var zeroNonce [chacha.NonceSize]byte

const (
	handshakeStatusEmpty       byte = 0
	handshakeStatusInitialized byte = 1
	handshakeStatusReceived    byte = 2
	handshakeStatusResponded   byte = 3
	handshakeStatusCompleted   byte = 4
	handshakeStatusFailed      byte = 5
	handshakeStatusSuccessful  byte = 6
)

// Handshake contains all the information exchanged between two nodes during
// both initialization and response. The state of a handshake struct must be
// identical for both nodes by the end of the protocol.
type Handshake struct {
	status byte // Stage at which the handshake is

	// for locking handshake while working on it
	sync.RWMutex

	Hash    crypt.Blake2sHash  // Chaining hash
	c, t, k [blake2s.Size]byte // Chaining keys

	keychain struct {
		sk ppk.PrivateKey // Handshake unique private key
		pk ppk.PublicKey  // Handshake unique public key
	}
}

// Reset sets all the fields of a handshake to zero.
func (h *Handshake) Reset() {
	h.Lock()
	crypt.ZeroBytes(h.Hash[:], h.c[:], h.k[:], h.t[:],
		h.keychain.sk[:], h.keychain.pk[:])
	h.status = handshakeStatusEmpty
	h.Unlock()
}

// Initialize will prepare the handshake for being sent as a message to a
// foreign node. The node will then respond and the handshake will proceed.
func (h *Handshake) Initialize(
	initializerSSK ppk.PrivateKey, receiverSPK ppk.PublicKey) *MsgHandshakeI {

	// lock handshake
	h.Lock()
	defer h.Unlock()

	// only initialize empty handshakes
	if h.status != handshakeStatusEmpty {
		return nil
	}

	// hash receiver public key
	crypt.Hash(&h.Hash, receiverSPK[:])

	// create unique keys for session
	h.keychain.sk, _ = ppk.NewPrivateKey()
	h.keychain.pk = h.keychain.sk.PublicKey()

	// derive keys
	hkdf.HKDF(h.Hash[:], h.keychain.pk[:], &h.c)

	// compute shared secret
	// between unique private key and receiver static public key
	ss := h.keychain.sk.SharedSecret(receiverSPK)

	// derive keys
	hkdf.HKDF(h.c[:], ss[:], &h.c, &h.k)

	// create message
	msg := new(MsgHandshakeI)
	msg.PlaintextUniquePublicKey = h.keychain.pk

	// ciphertext
	initializerSPK := initializerSSK.PublicKey()
	cipher, _ := chacha.New(h.k[:])
	cipher.Seal(msg.EncryptedStaticPublicKey[:0],
		zeroNonce[:], initializerSPK[:], h.Hash[:])

	// mark as initialized
	h.status = handshakeStatusInitialized

	return msg
}

// Receive is called with the current static private key and the incoming
// initialization handshake message from a peer. This will update the state
// of the handshake and authenticate the foreign peer and their key.
func (h *Handshake) Receive(
	msg *MsgHandshakeI, receiverSSK ppk.PrivateKey) (ppk.PublicKey, error) {

	// lock handshake
	h.Lock()
	defer h.Unlock()

	// hash receiver public key
	receiverSPK := receiverSSK.PublicKey()
	crypt.Hash(&h.Hash, receiverSPK[:])

	// derive keys
	hkdf.HKDF(h.Hash[:], msg.PlaintextUniquePublicKey[:], &h.c)

	// compute shared secret
	// between static private key and unique public key of initializer
	ss := receiverSSK.SharedSecret(msg.PlaintextUniquePublicKey)

	// derive keys
	hkdf.HKDF(h.c[:], ss[:], &h.c, &h.k)

	// return peer public static key
	var initializerSPK ppk.PublicKey

	// attempt to decrypt ciphertext and validate authenticity
	cipher, _ := chacha.New(h.k[:])
	_, err := cipher.Open(initializerSPK[:0],
		zeroNonce[:], msg.EncryptedStaticPublicKey[:], h.Hash[:])

	// if decryption fails, mark handshake as failed
	if err != nil {
		h.status = handshakeStatusFailed
		return initializerSPK, err
	}

	h.status = handshakeStatusReceived
	return initializerSPK, nil
}

// Respond is applied when a valid initialization message was received. This
// is the last step in establishing secure communication between the two nodes.
func (h *Handshake) Respond(
	initializerUPK, initializerSPK ppk.PublicKey) *MsgHandshakeR {

	// lock handshake
	h.Lock()
	defer h.Unlock()

	// only respond to received handshakes
	if h.status != handshakeStatusReceived {
		return nil
	}

	// create unique keys for session
	h.keychain.sk, _ = ppk.NewPrivateKey()
	h.keychain.pk = h.keychain.sk.PublicKey()

	// derive keys
	hkdf.HKDF(h.c[:], h.keychain.pk[:], &h.c)

	// compute hash
	crypt.Hash(&h.Hash, h.Hash[:], h.keychain.pk[:])

	// derive keys
	var ss [ppk.KeySize]byte
	ss = h.keychain.sk.SharedSecret(initializerUPK)
	hkdf.HKDF(h.c[:], ss[:], &h.c)
	ss = h.keychain.sk.SharedSecret(initializerSPK)
	hkdf.HKDF(h.c[:], ss[:], &h.c)

	// derive keys
	// TODO Change []byte{0} to pre-shared key optional similar to WireGuard
	hkdf.HKDF(h.c[:], []byte{0}, &h.c, &h.t, &h.k)

	// update hash
	crypt.Hash(&h.Hash, h.Hash[:], h.t[:])

	// create message
	msg := new(MsgHandshakeR)
	msg.PlaintextUniquePublicKey = h.keychain.pk

	// ciphertext
	cipher, _ := chacha.New(h.k[:])
	cipher.Seal(msg.EncryptedNothing[:0],
		zeroNonce[:], nil, h.Hash[:])

	// update hash
	crypt.Hash(&h.Hash, msg.EncryptedNothing[:])

	h.status = handshakeStatusResponded
	return msg
}

// Complete is the last step in the handshake protocol. The action is taken by
// the initializer.
func (h *Handshake) Complete(
	msg *MsgHandshakeR, initializerSSK ppk.PrivateKey) error {

	// lock handshake
	h.Lock()
	defer h.Unlock()

	// derive keys
	hkdf.HKDF(h.c[:], msg.PlaintextUniquePublicKey[:], &h.c)

	// compute hash
	crypt.Hash(&h.Hash, h.Hash[:], msg.PlaintextUniquePublicKey[:])

	// derive keys
	var ss [ppk.KeySize]byte
	ss = h.keychain.sk.SharedSecret(msg.PlaintextUniquePublicKey)
	hkdf.HKDF(h.c[:], ss[:], &h.c)
	ss = initializerSSK.SharedSecret(msg.PlaintextUniquePublicKey)
	hkdf.HKDF(h.c[:], ss[:], &h.c)

	// derive keys
	// TODO Change []byte{0} to pre-shared key optional similar to WireGuard
	hkdf.HKDF(h.c[:], []byte{0}, &h.c, &h.t, &h.k)

	// update hash
	crypt.Hash(&h.Hash, h.Hash[:], h.t[:])

	// ciphertext
	cipher, _ := chacha.New(h.k[:])
	_, err := cipher.Open(nil, zeroNonce[:], msg.EncryptedNothing[:], h.Hash[:])

	// if decryption fails, mark handshake as failed
	if err != nil {
		h.status = handshakeStatusFailed
		return err
	}

	// update hash
	crypt.Hash(&h.Hash, msg.EncryptedNothing[:])

	h.status = handshakeStatusCompleted
	return nil
}

// Finalize is called after the handshake protocol is either completed by an
// initializer or responded to by the receiver. In this step the transport keys
// are derived from the final state of the handshake and other future unused
// fields are zeroed.
func (h *Handshake) Finalize() (send, recv [ppk.KeySize]byte) {
	// lock handshake
	h.Lock()
	defer h.Unlock()

	// check for initializer or receiver mode
	// default case send, recv keys are all 0
	// Finalize works if and only if the handshake status is Successful
	switch h.status {
	// initializer
	case handshakeStatusCompleted:
		hkdf.HKDF(h.c[:], nil, &send, &recv)
		h.status = handshakeStatusSuccessful
	// receiver
	case handshakeStatusResponded:
		hkdf.HKDF(h.c[:], nil, &recv, &send)
		h.status = handshakeStatusSuccessful
	}

	// zero unused data
	crypt.ZeroBytes(h.k[:], h.c[:], h.keychain.pk[:], h.keychain.sk[:], h.t[:])

	return
}
