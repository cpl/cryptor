package proto

import (
	"bytes"
	"encoding/gob"

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
	// 80 = PlaintextUniquePublicKey (32) + EncryptedStaticPublicKey (48)
	MsgSizeHandshakeI = 80

	// MsgSizeHandshakeR - Handshake Response Size
	// 0 = ?
	// TODO Define handshake response structure and compute size
	MsgSizeHandshakeR = 0
)

// MsgHandshakeI is the first message used by an initiator on the
// Cryptor Network to contact and establish secure communication with a peer.
type MsgHandshakeI struct {
	PlaintextUniquePublicKey ppk.PublicKey
	EncryptedStaticPublicKey [48]byte
}

// MarshalBinary returns the content of the message as a single byte array.
// TODO Improve binary marshaling of messages in general
func (m *MsgHandshakeI) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m)
	return buf.Bytes(), err
}
