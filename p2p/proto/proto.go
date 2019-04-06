/*
Package proto defines the Cryptor P2P Network base protocol. It contains the
algorithms and structures used in communication.
*/
package proto // import "cpl.li/go/cryptor/p2p/proto"

import (
	"golang.org/x/crypto/blake2s"
	chacha "golang.org/x/crypto/chacha20poly1305"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/hkdf"
	"cpl.li/go/cryptor/crypt/ppk"
)

// NewMsgHandshakeI creates a new handshake initiation message for
// the peer with the given public key.
func NewMsgHandshakeI(pkr ppk.PublicKey) *MsgHandshakeI {
	// data
	msg := new(MsgHandshakeI)
	var h crypt.Blake2sHash
	var c, k [blake2s.Size]byte
	var zeroNonce [chacha.NonceSize]byte

	// ! DEBUG
	// zero data at the end
	defer crypt.ZeroBytes(h[:], c[:], k[:])

	// hash receiver public key
	crypt.Hash(&h, pkr[:])

	// ! DEBUG
	// generate unique keypair
	ski, _ := ppk.NewPrivateKey()
	pki := ski.PublicKey()

	// derive key
	hkdf.HKDF(h[:], pki[:], &c)

	// compute shared secret from unique and static
	ss := ski.SharedSecret(pkr)

	// derive key
	hkdf.HKDF(c[:], ss[:], &c, &k)

	// ciphertext
	cipher, _ := chacha.New(k[:])
	cipher.Seal(msg.EncryptedStaticPublicKey[:], zeroNonce[:], pkr[:], h[:])

	// set message fields
	msg.PlaintextUniquePublicKey = pki

	return msg
}

// Validate checks weather the incoming message is valid and intended for the
// local public key.
// TODO Change Validation flow to integrate more nicely with the Node
// ! DEBUG
func (msg *MsgHandshakeI) Validate(sk ppk.PrivateKey) (ppk.PublicKey, bool) {
	// data
	pk := sk.PublicKey()
	var pki ppk.PublicKey
	var h crypt.Blake2sHash
	var c, k [blake2s.Size]byte
	var zeroNonce [chacha.NonceSize]byte

	// ! DEBUG
	// zero data at the end
	defer crypt.ZeroBytes(h[:], c[:], k[:])

	// hash public key
	crypt.Hash(&h, pk[:])

	// derive key
	hkdf.HKDF(h[:], msg.PlaintextUniquePublicKey[:], &c)

	// compute shared secret from unique and static
	ss := sk.SharedSecret(msg.PlaintextUniquePublicKey)

	// derive key
	hkdf.HKDF(c[:], ss[:], &c, &k)

	// ciphertext
	cipher, _ := chacha.New(k[:])
	_, err :=
		cipher.Open(pki[:], zeroNonce[:], msg.EncryptedStaticPublicKey[:], h[:])
	if err != nil {
		return pki, false
	}

	return pki, true
}
