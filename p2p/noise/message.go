package noise

import (
	"encoding/binary"
	"errors"

	"cpl.li/go/cryptor/crypt/ppk"
)

const (
	// SizeMessageInitializer is the size of a handshake message sent by
	// initializer to the responder.
	SizeMessageInitializer = sizeID + ppk.KeySize + sizeEncPub

	// SizeMessageResponder is the size of the response message from the
	// responder to the initializer.
	SizeMessageResponder = sizeID + ppk.KeySize + sizeEncNth
)

// encryption sizes
const (
	sizeEncPub = 48 // encrypted size of static public key
	sizeEncNth = 16 // encrypted size of nothing (nil)
	sizeID     = 8  // size of an int64 (64/8 = 8 bytes)
)

// MessageInitializer encapsulates the data which is sent by the initializer to
// the responder.
type MessageInitializer struct {
	ID                                  uint64
	PlaintextUniquePublic               ppk.PublicKey
	EncryptedInitializerStaticPublicKey [sizeEncPub]byte
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (msg *MessageInitializer) MarshalBinary() ([]byte, error) {
	// create output buffer
	out := make([]byte, SizeMessageInitializer)

	// write fields
	binary.LittleEndian.PutUint64(out, msg.ID)
	copy(out[sizeID:], msg.PlaintextUniquePublic[:])
	copy(out[sizeID+ppk.KeySize:], msg.EncryptedInitializerStaticPublicKey[:])

	return out, nil
}

// UnmarshalBinary can unmarshal the output of MarshalBinary back into itself.
func (msg *MessageInitializer) UnmarshalBinary(data []byte) error {
	// check data size
	if len(data) != SizeMessageInitializer {
		return errors.New("invalid message size")
	}

	// unpack fields
	msg.ID = binary.LittleEndian.Uint64(data)
	copy(msg.PlaintextUniquePublic[:], data[sizeID:sizeID+ppk.KeySize])
	copy(msg.EncryptedInitializerStaticPublicKey[:], data[sizeID+ppk.KeySize:])

	return nil
}

// MessageResponder encapsulates the data which is sent by the responder to
// the initializer.
type MessageResponder struct {
	ID                    uint64
	PlaintextUniquePublic ppk.PublicKey
	EncryptedNothing      [sizeEncNth]byte
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (msg *MessageResponder) MarshalBinary() ([]byte, error) {
	// create output buffer
	out := make([]byte, SizeMessageResponder)

	// write fields
	binary.LittleEndian.PutUint64(out, msg.ID)
	copy(out[sizeID:], msg.PlaintextUniquePublic[:])
	copy(out[sizeID+ppk.KeySize:], msg.EncryptedNothing[:])

	return out, nil
}

// UnmarshalBinary can unmarshal the output of MarshalBinary back into itself.
func (msg *MessageResponder) UnmarshalBinary(data []byte) error {
	// check data size
	if len(data) != SizeMessageResponder {
		return errors.New("invalid message size")
	}

	// unpack fields
	msg.ID = binary.LittleEndian.Uint64(data)
	copy(msg.PlaintextUniquePublic[:], data[sizeID:sizeID+ppk.KeySize])
	copy(msg.EncryptedNothing[:], data[sizeID+ppk.KeySize:])

	return nil
}
