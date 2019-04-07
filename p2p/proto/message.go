package proto

import (
	"cpl.li/go/cryptor/crypt/ppk"
)

const (
	// MsgTypeHandshakeI - Handshake Initializer
	MsgTypeHandshakeI byte = 0

	// MsgTypeHandshakeR - Handshake Response
	MsgTypeHandshakeR byte = 1
)

const (
	// MsgSizeHandshakeI - Handshake Initializer Size
	// 81 = PlaintextUniquePublicKey (32) + EncryptedStaticPublicKey (48)
	MsgSizeHandshakeI = 80

	// MsgSizeHandshakeR - Handshake Response Size
	// 48 = PlaintextUniquePublicKey (32) + EncryptedNothing (16)
	MsgSizeHandshakeR = 48
)

// MsgHandshakeI is the first message used by an initiator on the
// Cryptor Network to contact and establish secure communication with a peer.
type MsgHandshakeI struct {
	PlaintextUniquePublicKey ppk.PublicKey
	EncryptedStaticPublicKey [48]byte
}

// MsgHandshakeR is the response created by the receiver to the initializer
// after the first handshake was successfully validated.
type MsgHandshakeR struct {
	PlaintextUniquePublicKey ppk.PublicKey
	EncryptedNothing         [16]byte
}
