package net_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/blake2s"
	chacha "golang.org/x/crypto/chacha20poly1305"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/hkdf"
	"cpl.li/go/cryptor/crypt/ppk"
)

type keypair struct {
	privateKey ppk.PrivateKey
	publicKey  ppk.PublicKey
}

func newKeypair() *keypair {
	kp := new(keypair)
	kp.privateKey, _ = ppk.NewPrivateKey()
	kp.publicKey = kp.privateKey.PublicKey()
	return kp
}

func TestProtocolCrypto(t *testing.T) {
	// TODO Remove this test & start developing the custom Noise implementation

	// in the following scenario we have 2 keypairs (or two bare-bone
	// representations of two nodes)
	// keypair 0 is aware of the static public key of keypair node 1
	// keypair 0 is the initiator and keypair 1 is the responder
	kp0 := newKeypair()
	kp1 := newKeypair()

	// keypair 0 generates a new set of keys, uniquely used for this session
	// with keypair 1
	kp0u := newKeypair()

	// the protocol contains a constant identifier
	identifier := []byte("cryptor-cpl")

	// the identifier is hashed and used as a key in HKDF
	var h crypt.Blake2sHash
	crypt.Hash(&h, identifier, kp1.publicKey[:])
	fmt.Println(h.ToHex())

	// start key for AEAD and chaining
	var c, k [blake2s.Size]byte

	// generate chaining key from static hashed identifier and unique pub key
	hkdf.HKDF(h[:], kp0u.publicKey[:], &c)

	// generate shared secret for unique key and foreign pub key
	ss0 := kp0u.privateKey.SharedSecret(kp1.publicKey)
	fmt.Println(hex.EncodeToString(ss0[:]))
	// chain key and generate key for AEAD
	hkdf.HKDF(c[:], ss0[:], &c, &k)

	// start AEAD encryption using k and 0 nonce
	var zeroNonce [chacha.NonceSize]byte
	cipher, _ := chacha.New(k[:])

	// generate encrypted text with kp0 public key
	ciphertext := cipher.Seal(nil, zeroNonce[:], kp0.publicKey[:], h[:])
	fmt.Println(kp0.publicKey.ToHex())

	// update hash
	crypt.Hash(&h, h[:], ciphertext)
	fmt.Println(h.ToHex())

	// generate shared secret from static keys
	ss1 := kp0.privateKey.SharedSecret(kp1.publicKey)
	fmt.Println(hex.EncodeToString(ss1[:]))

	// update chaining keys
	hkdf.HKDF(c[:], ss1[:], &c, &k)

	// new cipher
	timestampbytes, _ := time.Now().MarshalBinary()
	cipher, _ = chacha.New(k[:])
	timestamp := cipher.Seal(nil, zeroNonce[:], timestampbytes, h[:])

	// update hash
	crypt.Hash(&h, h[:], timestamp)
	fmt.Println(h.ToHex())

	// by this point, the initial message is complete:
	// - kp0u.public
	// - ciphertext
	// - timestamp
	fmt.Println("----------------")

	// it's time for the receiving node to apply the same protocol and
	// get to the same state as the init node
	var hr crypt.Blake2sHash
	crypt.Hash(&hr, identifier, kp1.publicKey[:])
	fmt.Println(hr.ToHex())
	var cr, kr [blake2s.Size]byte
	hkdf.HKDF(hr[:], kp0u.publicKey[:], &cr) // kp0u.public is in message
	ssr0 := kp1.privateKey.SharedSecret(kp0u.publicKey)
	fmt.Println(hex.EncodeToString(ssr0[:]))
	hkdf.HKDF(cr[:], ssr0[:], &cr, &kr)
	cipher, _ = chacha.New(kr[:])
	plaintext, err := cipher.Open(nil, zeroNonce[:], ciphertext, hr[:])
	if err != nil {
		t.Fatal(err)
	}

	// a debug check
	if !bytes.Equal(plaintext, kp0.publicKey[:]) {
		t.Fatalf("plaintext does not match public key\n")
	}

	// at this point we have access to kp0.public
	fmt.Println(hex.EncodeToString(plaintext))

	// hash encrypted received message
	crypt.Hash(&hr, hr[:], ciphertext)
	fmt.Println(hr.ToHex())

	// compute shared secret for static keys
	ssr1 := kp1.privateKey.SharedSecret(kp0.publicKey)
	fmt.Println(hex.EncodeToString(ssr1[:]))

	// update chaining keys
	hkdf.HKDF(cr[:], ssr1[:], &cr, &kr)

	cipher, _ = chacha.New(kr[:])
	plaintimestamp, err := cipher.Open(nil, zeroNonce[:], timestamp, hr[:])
	if err != nil {
		t.Fatal(err)
	}

	// debug check
	if !bytes.Equal(plaintimestamp, timestampbytes) {
		t.Fatalf("failed to verify timestamp")
	}
}
