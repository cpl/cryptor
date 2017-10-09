package p2p

import (
	"net"

	"github.com/thee-engineer/cryptor/network"
)

// Peer represents a foreign node connected to the local node.
type Peer struct {
	addr *net.UDPAddr
}

// NewPeer creates a peer object given an IP:PORT pair. This is for testing and
// debugging purposes. Not to be used in production.
func NewPeer(ip string, port int) *Peer {
	return &Peer{
		addr: network.IPPToUDP(ip, port),
	}
}
