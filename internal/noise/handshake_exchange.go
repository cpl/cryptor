package noise

import (
	"cpl.li/go/cryptor/internal/crypt"
	"cpl.li/go/cryptor/internal/crypt/hashing"
	"cpl.li/go/cryptor/internal/crypt/hkdf"
	"cpl.li/go/cryptor/internal/crypt/ppk"

	chacha "golang.org/x/crypto/chacha20poly1305"
)

// InitializeSender ...
func (hs *Handshake) InitializeSender(rPub *ppk.PublicKey) error {
	if hs.state != handshakeStateEmpty {
		return ErrBadHandshakeState
	}

	if err := hs.keygen(); err != nil {
		return err
	}

	hashing.Hash(&hs.hash, rPub[:])
	hkdf.HKDF(hs.hash[:], hs.tempKeys.public[:], &hs.c)

	var ss [ppk.KeySize]byte

	hs.tempKeys.secret.SharedSecret(rPub, &ss)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c, &hs.k)

	hs.role = handshakeRoleSender
	hs.state = handshakeStateInitialized

	return nil
}

// InitializeRecipient ...
func (hs *Handshake) InitializeRecipient(rSec *ppk.PrivateKey, cPubTmp *ppk.PublicKey) error {
	var rPub ppk.PublicKey
	rSec.PublicKey(&rPub)

	hashing.Hash(&hs.hash, rPub[:])
	hkdf.HKDF(hs.hash[:], cPubTmp[:], &hs.c)

	var ss [ppk.KeySize]byte

	rSec.SharedSecret(cPubTmp, &ss)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c, &hs.k)

	if err := hs.keygen(); err != nil {
		return err
	}

	hs.role = handshakeRoleRecipient
	hs.state = handshakeStateInitialized

	return nil
}

// Exchange ...
func (hs *Handshake) Exchange(cPub *ppk.PublicKey, cPubEnc *EncryptedKey) error {
	if hs.state != handshakeStateInitialized {
		return ErrBadHandshakeState
	}

	switch hs.role {
	case handshakeRoleSender:
		cipher, _ := chacha.New(hs.k[:])
		cipher.Seal(cPubEnc[:0], zeroNonce[:], cPub[:], hs.hash[:])
	case handshakeRoleRecipient:
		cipher, _ := chacha.New(hs.k[:])
		_, err := cipher.Open(cPub[:0], zeroNonce[:], cPubEnc[:], hs.hash[:])
		if err != nil {
			return err
		}
	default:
		return ErrBadHandshakeRole
	}

	hs.state = handshakeStateExchanged

	return nil
}

// PrepareRecipientResponse ...
func (hs *Handshake) PrepareRecipientResponse(sPubTmp, sPub *ppk.PublicKey, enc *EncryptedNothing) error {
	if hs.role != handshakeRoleRecipient {
		return ErrBadHandshakeRole
	}

	if hs.state != handshakeStateExchanged {
		return ErrBadHandshakeState
	}

	hkdf.HKDF(hs.c[:], hs.tempKeys.public[:], &hs.c)
	hashing.Hash(&hs.hash, hs.hash[:], hs.tempKeys.public[:])

	var ss [ppk.KeySize]byte

	hs.tempKeys.secret.SharedSecret(sPubTmp, &ss)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c)
	hs.tempKeys.secret.SharedSecret(sPub, &ss)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c)

	hkdf.HKDF(hs.c[:], hs.presharedKey[:], &hs.c, &hs.t, &hs.k)
	hashing.Hash(&hs.hash, hs.hash[:], hs.t[:])

	cipher, _ := chacha.New(hs.k[:])
	cipher.Seal(enc[:0], zeroNonce[:], nil, hs.hash[:])

	hashing.Hash(&hs.hash, enc[:])

	hs.state = handshakeStateFinal

	return nil
}

// ConsumeRecipientResponse ...
func (hs *Handshake) ConsumeRecipientResponse(sSec *ppk.PrivateKey, rPubTmp *ppk.PublicKey, enc *EncryptedNothing) error {
	if hs.role != handshakeRoleSender {
		return ErrBadHandshakeRole
	}

	if hs.state != handshakeStateExchanged {
		return ErrBadHandshakeState
	}

	var (
		snapshotHash = hs.hash
		snapshotC    = hs.c
		snapshotK    = hs.k
		snapshotT    = hs.t
	)

	hkdf.HKDF(hs.c[:], rPubTmp[:], &hs.c)
	hashing.Hash(&hs.hash, hs.hash[:], rPubTmp[:])

	var ss [ppk.KeySize]byte

	hs.tempKeys.secret.SharedSecret(rPubTmp, &ss)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c)
	sSec.SharedSecret(rPubTmp, &ss)
	hkdf.HKDF(hs.c[:], ss[:], &hs.c)

	hkdf.HKDF(hs.c[:], hs.presharedKey[:], &hs.c, &hs.t, &hs.k)
	hashing.Hash(&hs.hash, hs.hash[:], hs.t[:])

	cipher, _ := chacha.New(hs.k[:])
	_, err := cipher.Open(nil, zeroNonce[:], enc[:], hs.hash[:])
	if err != nil {
		hs.hash = snapshotHash
		hs.c = snapshotC
		hs.k = snapshotK
		hs.t = snapshotT

		return err
	}

	hashing.Hash(&hs.hash, enc[:])

	hs.state = handshakeStateFinal

	return nil
}

// Finalize ...
func (hs *Handshake) Finalize(send, recv *[ppk.KeySize]byte) error {
	if hs.state != handshakeStateFinal {
		return ErrBadHandshakeState
	}

	switch hs.role {
	case handshakeRoleSender:
		hkdf.HKDF(hs.c[:], nil, send, recv)
	case handshakeRoleRecipient:
		hkdf.HKDF(hs.c[:], nil, recv, send)
	default:
		return ErrBadHandshakeRole
	}

	hs.state = handshakeStateComplete

	crypt.ZeroBytes(hs.k[:], hs.c[:], hs.t[:],
		hs.tempKeys.public[:], hs.tempKeys.secret[:])

	return nil
}
