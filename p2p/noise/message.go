package noise

import (
	"errors"

	"cpl.li/go/cryptor/crypt/ppk"
)

// MessageInitializer encapsulates the data which is sent by the initializer to
// the responder.
type MessageInitializer struct {
	PlaintextUniquePublic               ppk.PublicKey
	EncryptedInitializerStaticPublicKey [sizeEncPub]byte
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (msg *MessageInitializer) MarshalBinary() ([]byte, error) {
	return append(
		msg.PlaintextUniquePublic[:],
		msg.EncryptedInitializerStaticPublicKey[:]...), nil
}

// UnmarshalBinary can unmarshal the output of MarshalBinary back into itself.
func (msg *MessageInitializer) UnmarshalBinary(data []byte) error {
	// check data size
	if len(data) != HandshakeSizeInitializer {
		return errors.New("invalid message size")
	}

	copy(msg.PlaintextUniquePublic[:], data[:ppk.KeySize])
	copy(msg.EncryptedInitializerStaticPublicKey[:], data[ppk.KeySize:])

	return nil
}

// MessageResponder encapsulates the data which is sent by the responder to
// the initializer.
type MessageResponder struct {
	PlaintextUniquePublic ppk.PublicKey
	EncryptedNothing      [sizeEncNth]byte
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (msg *MessageResponder) MarshalBinary() ([]byte, error) {
	return append(
		msg.PlaintextUniquePublic[:],
		msg.EncryptedNothing[:]...), nil
}

// UnmarshalBinary can unmarshal the output of MarshalBinary back into itself.
func (msg *MessageResponder) UnmarshalBinary(data []byte) error {
	// check data size
	if len(data) != HandshakeSizeResponder {
		return errors.New("invalid message size")
	}

	copy(msg.PlaintextUniquePublic[:], data[:ppk.KeySize])
	copy(msg.EncryptedNothing[:], data[ppk.KeySize:])

	return nil
}
