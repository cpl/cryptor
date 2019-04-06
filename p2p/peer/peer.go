/*
Package peer offers an abstraction layer for foreign nodes.
*/
package peer // import "cpl.li/go/cryptor/p2p/peer"

import (
	"net"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p"
)

// Peer represents a foreign machine/node on the Cryptor network and all the
// required information to establish a secure connection.
type Peer struct {
	addr *net.UDPAddr // address of foreign node

	staticPublicKey    ppk.PublicKey  // static public key (identifier)
	uniquePrivateKey   ppk.PrivateKey // connection unique private key
	uniqueSharedSecret [32]byte       // connection unique shared secret

	isConfirmed bool // connection with peer was confirmed with a handshake
}

// New creates a new peer with the given public key.
func New(pk ppk.PublicKey) *Peer {
	p := new(Peer)

	p.staticPublicKey = pk

	return p
}

// PublicKey returns the static public key of the peer.
func (p *Peer) PublicKey() ppk.PublicKey {
	return p.PublicKey()
}

// Addr returns the remote network address of the peer as ip:port.
func (p *Peer) Addr() string {
	return p.addr.String()
}

// SetAddr sets the remote network address of the peer.
func (p *Peer) SetAddr(addr string) (err error) {
	p.addr, err = net.ResolveUDPAddr(p2p.Network, addr)
	return err
}
