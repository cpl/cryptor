/*
Package peer implements the representation of a foreign Cryptor node.
*/
package peer // import "cpl.li/go/cryptor/p2p/peer"

import (
	"net"
	"sync"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p"
	"cpl.li/go/cryptor/p2p/noise"
)

// Peer represents a foreign machine/node on the Cryptor network and all the
// required information to establish a secure connection.
type Peer struct {
	id              uint64        // id used in network communication
	addr            *net.UDPAddr  // address of foreign node
	staticPublicKey ppk.PublicKey // static public key (identifier)

	// lock for operating on the peer
	sync.RWMutex

	HasHandshake bool
	Handshake    *noise.Handshake // the current handshake status with this peer

	// transport keys generated from the finalization of the handshake
	keys struct {
		send [ppk.KeySize]byte
		recv [ppk.KeySize]byte
	}
}

// NewPeer creates a new peer with the given public key and optional address.
func NewPeer(pk ppk.PublicKey, addr string) *Peer {
	// create peer
	p := new(Peer)

	// set public key
	p.staticPublicKey = pk

	// set handshake as nil
	p.Handshake = nil

	// set address if any is given
	if addr != "" {
		p.addr, _ = net.ResolveUDPAddr(p2p.Network, addr)
	}

	// assign random ID
	p.id = crypt.RandomUint64()

	return p
}

// Addr returns the string of the address.
func (p *Peer) Addr() string {
	return p.addr.String()
}

// AddrUDP returns the UDP address of the peer.
func (p *Peer) AddrUDP() *net.UDPAddr {
	return p.addr
}

// PublicKey returns the peer known static public key.
func (p *Peer) PublicKey() ppk.PublicKey {
	return p.staticPublicKey
}
