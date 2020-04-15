package noise

import (
	"cpl.li/go/cryptor/internal/crypt/hashing"
	"cpl.li/go/cryptor/internal/crypt/ppk"
)

type handshakeRole byte

const (
	handshakeRoleEmpty handshakeRole = iota
	handshakeRoleSender
	handshakeRoleRecipient
)

type handshakeState byte

const (
	handshakeStateEmpty handshakeState = iota
	handshakeStateInitialized
	handshakeStateExchanged
	handshakeStateFinal
	handshakeStateComplete
)

// Handshake ...
type Handshake struct {
	state handshakeState
	role  handshakeRole

	hash    hashing.HashSum
	c, t, k [hashing.HashSize]byte

	presharedKey [ppk.KeySize]byte

	tempKeys struct {
		secret ppk.PrivateKey
		public ppk.PublicKey
	}
}

// PublicKey ...
func (hs *Handshake) PublicKey() ppk.PublicKey {
	return hs.tempKeys.public
}

func (hs *Handshake) keygen() (err error) {
	if err := ppk.NewPrivateKey(&hs.tempKeys.secret); err != nil {
		return err
	}
	if err := hs.tempKeys.secret.PublicKey(&hs.tempKeys.public); err != nil {
		return err
	}

	return nil
}
