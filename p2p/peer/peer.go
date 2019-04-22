/*
Package peer implements the representation of a foreign Cryptor node.
*/
package peer // import "cpl.li/go/cryptor/p2p/peer"

import (
	"net"
	"sync"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p"
	"cpl.li/go/cryptor/p2p/noise"
)

// Peer represents a foreign machine/node on the Cryptor network and all the
// required information to establish a secure connection.
type Peer struct {
	ID   uint64       // id used in network communication
	addr *net.UDPAddr // address of foreign node

	// lock for operating on the peer
	sync.RWMutex

	Handshake *noise.Handshake // the current handshake with this peer

	// transport keys generated from the finalization of the handshake
	keys struct {
		sync.RWMutex

		staticPublicKey ppk.PublicKey     // static public key (identifier)
		send            [ppk.KeySize]byte // for sending transport messages
		recv            [ppk.KeySize]byte // for receiving transport messages
	}
}

// NewPeer creates a new peer with the given public key and optional address.
func NewPeer(pk ppk.PublicKey, addr string) *Peer {
	// create peer
	p := new(Peer)

	// set public key
	p.keys.staticPublicKey = pk

	// set handshake as nil
	p.Handshake = nil

	// set address if any is given
	if addr != "" {
		p.addr, _ = net.ResolveUDPAddr(p2p.Network, addr)
	}

	return p
}

// Addr returns the string of the address.
func (p *Peer) Addr() string {
	// lock peer
	p.RLock()
	defer p.RUnlock()

	return p.addr.String()
}

// AddrUDP returns the UDP address of the peer.
func (p *Peer) AddrUDP() *net.UDPAddr {
	// lock peer
	p.RLock()
	defer p.RUnlock()

	return p.addr
}

// PublicKey returns the peer known static public key.
func (p *Peer) PublicKey() ppk.PublicKey {
	// lock keys
	p.keys.RLock()
	defer p.keys.RUnlock()

	// return
	return p.keys.staticPublicKey
}

// SetAddr sets the remote network address of the peer.
func (p *Peer) SetAddr(addr string) error {
	// lock peer
	p.Lock()
	defer p.Unlock()

	// set to nil
	if addr == "" {
		p.addr = nil
		return nil
	}

	// resolve addr
	newaddr, err := net.ResolveUDPAddr(p2p.Network, addr)
	if err != nil {
		return err
	}

	// set addr
	p.addr = newaddr

	return nil
}

// SetTransportKeys sets the keys used for encryption and decryption of
// outgoing and incoming transport messages.
func (p *Peer) SetTransportKeys(send, recv [ppk.KeySize]byte) {
	// lock keys
	p.keys.Lock()
	defer p.keys.Unlock()

	// set keys
	p.keys.send = send
	p.keys.recv = recv

	return
}
