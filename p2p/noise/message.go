package noise

import (
	"cpl.li/go/cryptor/crypt/ppk"
)

// MessageInitializer encapsulates the data which is sent by the initializer to
// the responder.
type MessageInitializer struct {
	PlaintextUniquePublic               ppk.PublicKey
	EncryptedInitializerStaticPublicKey [sizeEncPub]byte
}

// MessageResponder encapsulates the data which is sent by the responder to
// the initializer.
type MessageResponder struct {
	PlaintextUniquePublic ppk.PublicKey
	EncryptedNothing      [sizeEncNth]byte
}

// Message is an abstraction for MessageInitializer and MessageResponder.
type Message interface {
	Bytes() []byte
	PublicKey() []byte
	Payload() []byte
}
