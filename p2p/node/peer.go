package node

import (
	"net"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p"
	"cpl.li/go/cryptor/p2p/proto"
)

// Peer represents a foreign machine/node on the Cryptor network and all the
// required information to establish a secure connection.
type Peer struct {
	addr            *net.UDPAddr  // address of foreign node
	staticPublicKey ppk.PublicKey // static public key (identifier)

	Handshake proto.Handshake
}

// NewPeer creates a new peer with the given public key and optional address.
func NewPeer(pk ppk.PublicKey, addr string) *Peer {
	// create peer
	p := new(Peer)

	// set public key
	p.staticPublicKey = pk

	// set address if any is given
	if addr != "" {
		p.addr, _ = net.ResolveUDPAddr(p2p.Network, addr)
	}

	return p
}
